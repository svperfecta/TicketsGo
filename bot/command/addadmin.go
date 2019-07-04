package command

import (
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/database"
	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
)

type AddAdminCommand struct {
}

func (AddAdminCommand) Name() string {
	return "addadmin"
}

func (AddAdminCommand) Description() string {
	return "Grants a user admin privileges"
}

func (AddAdminCommand) Aliases() []string {
	return []string{}
}

func (AddAdminCommand) PermissionLevel() utils.PermissionLevel {
	return utils.Admin
}

func (AddAdminCommand) Execute(ctx CommandContext) {
	if len(ctx.Message.Mentions) == 0 {
		ctx.SendEmbed(utils.Red, "Error", "You need to mention a user to grant admin privileges to")
		ctx.ReactWithCross()
		return
	}

	var overwrites []*discordgo.PermissionOverwrite
	ch, err := ctx.Session.State.Channel(ctx.Channel); if err != nil {
		ch, err = ctx.Session.Channel(ctx.Channel); if err != nil {
			log.Error(err.Error())
			return
		}
	}

	overwrites = ch.PermissionOverwrites

	for _, mention := range ctx.Message.Mentions {
		go database.AddAdmin(ctx.Guild, mention.ID)

		overwrites = append(overwrites, &discordgo.PermissionOverwrite{
			ID: mention.ID,
			Type: "member",
			Allow: utils.SumPermissions(utils.ViewChannel, utils.SendMessages, utils.AddReactions, utils.AttachFiles, utils.ReadMessageHistory, utils.EmbedLinks),
			Deny: 0,
		})
	}

	data := discordgo.ChannelEdit{
		PermissionOverwrites: overwrites,
	}

	if _, err = ctx.Session.ChannelEditComplex(ctx.Guild, &data); err != nil {
		ctx.ReactWithCross()
		log.Error(err.Error())
		return
	}

	ctx.ReactWithCheck()
}

func (AddAdminCommand) Parent() interface{} {
	return nil
}

func (AddAdminCommand) Children() []Command {
	return make([]Command, 0)
}

func (AddAdminCommand) PremiumOnly() bool {
	return false
}

func (AddAdminCommand) AdminOnly() bool {
	return false
}

func (AddAdminCommand) HelperOnly() bool {
	return false
}
