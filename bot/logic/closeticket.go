package logic

import (
	"fmt"
	"github.com/TicketsBot/TicketsGo/bot/archive"
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/database"
	"github.com/TicketsBot/TicketsGo/sentry"
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/rest"
	"strconv"
	"strings"
)

func CloseTicket(s *gateway.Shard, guildId, channelId, messageId uint64, member member.Member, args []string, fromReaction, isPremium bool) {
	reference := message.MessageReference{
		MessageId: messageId,
		ChannelId: channelId,
		GuildId:   guildId,
	}

	// Get ticket struct
	ticket, err := database.Client.Tickets.GetByChannel(channelId)
	if err != nil {
		sentry.Error(err)
		return
	}

	isTicket := ticket.GuildId != 0

	// Verify this is a ticket or modmail channel
	// Cannot happen if fromReaction
	if !isTicket {
		// check whether this is a modmail channel
		var isModmail bool
		{
			modmailSession, err := database.Client.ModmailSession.GetByChannel(channelId)
			if err != nil {
				sentry.Error(err)
				return
			}

			isModmail = modmailSession.GuildId != 0
		}
		if isModmail {
			return
		}

		utils.ReactWithCross(s, reference)
		utils.SendEmbed(s, channelId, utils.Red, "Error", "This is not a ticket channel", nil, 30, isPremium)

		return
	}

	// Create reason
	var reason string
	silentClose := false
	for _, arg := range args {
		if arg == "--silent" {
			silentClose = true
		} else {
			reason += fmt.Sprintf("%s ", arg)
		}
	}
	reason = strings.TrimSuffix(reason, " ")

	// Check the user is permitted to close the ticket
	permissionLevelChan := make(chan utils.PermissionLevel)
	go utils.GetPermissionLevel(s, member, guildId, permissionLevelChan)
	permissionLevel := <-permissionLevelChan

	usersCanClose, err := database.Client.UsersCanClose.Get(guildId); if err != nil {
		sentry.Error(err)
	}

	if (permissionLevel == utils.Everyone && ticket.UserId != member.User.Id) || (permissionLevel == utils.Everyone && !usersCanClose) {
		if !fromReaction {
			utils.ReactWithCross(s, reference)
			utils.SendEmbed(s, channelId, utils.Red, "Error", "You are not permitted to close this ticket", nil, 30, isPremium)
		}
		return
	}

	// TODO: Re-add permission check
	/*if !permission.HasPermissions(s, guildId, s.SelfId(), permission.ManageChannels) {
		ctx.ReactWithCross()
		ctx.SendEmbed(utils.Red, "Error", "I do not have permission to delete this channel")
		return
	}*/

	if !fromReaction {
		utils.ReactWithCheck(s, reference)
	}

	// Archive
	msgs := make([]message.Message, 0)

	lastId := uint64(0)
	count := -1
	for count != 0 {
		array, err := s.GetChannelMessages(channelId, rest.GetChannelMessagesData{
			Before: lastId,
			Limit:  100,
		})

		count = len(array)
		if err != nil {
			count = 0
			sentry.Error(err)
		}

		if count > 0 {
			lastId = array[len(array)-1].Id

			msgs = append(msgs, array...)
		}
	}

	// Reverse messages
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	if err := archive.ArchiverClient.Store(msgs, guildId, ticket.Id, isPremium); err != nil {
		sentry.Error(err)
	}

	// Set ticket state as closed and delete channel
	if err :=  database.Client.Tickets.Close(ticket.Id, guildId); err != nil {
		sentry.Error(err)
	}

	if _, err := s.DeleteChannel(channelId); err != nil {
		sentry.Error(err)
	}

	// Send logs to archive channel
	archiveChannelId, err := database.Client.ArchiveChannel.Get(guildId); if err != nil {
		sentry.Error(err)
	}

	var channelExists bool
	if archiveChannelId != 0 {
		if _, err := s.GetChannel(archiveChannelId); err == nil {
			channelExists = true
		}
	}

	// Save space - delete the webhook
	go database.Client.Webhooks.Delete(guildId, ticket.Id)

	if channelExists {
		embed := embed.NewEmbed().
			SetTitle("Ticket Closed").
			SetColor(int(utils.Green)).
			AddField("Ticket ID", strconv.Itoa(ticket.Id), true).
			AddField("Closed By", member.User.Mention(), true).
			AddField("Archive", fmt.Sprintf("[Click here](https://panel.ticketsbot.net/manage/%d/logs/view/%d)", guildId, ticket.Id), true)

		if reason == "" {
			embed.AddField("Reason", "No reason specified", false)
		} else {
			embed.AddField("Reason", reason, false)
		}

		if _, err := s.CreateMessageEmbed(archiveChannelId, embed); err != nil {
			sentry.Error(err)
		}

		// Notify user and send logs in DMs
		if !silentClose {
			dmChannel, err := s.CreateDM(ticket.UserId)

			// Only send the msg if we could create the channel
			if err == nil {
				if _, err := s.CreateMessageEmbed(dmChannel.Id, embed); err != nil {
					sentry.Error(err)
				}
			}
		}
	}
}
