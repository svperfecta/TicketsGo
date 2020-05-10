package messagequeue

import (
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/cache"
	dbclient "github.com/TicketsBot/TicketsGo/database"
	"github.com/TicketsBot/TicketsGo/sentry"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/gateway"
	"github.com/rxdn/gdl/objects/channel/embed"
)

func ListenPanelCreations(shardManager *gateway.ShardManager) {
	creations := make(chan database.Panel)
	go cache.Client.ListenPanelCreate(creations)

	for panel := range creations {
		// Get shard
		shard := shardManager.ShardForGuild(panel.GuildId); if shard == nil {
			continue
		}

		errorContext := sentry.ErrorContext{
			Guild:       panel.GuildId,
			Channel:     panel.ChannelId,
			Shard:       shard.ShardId,
		}

		// Create embed object
		embed := embed.NewEmbed()

		// Get whether guild is premium
		isPremiumChan := make(chan bool)
		go utils.IsPremiumGuild(shard, panel.GuildId, isPremiumChan)
		isPremium := <-isPremiumChan

		if !isPremium {
			embed.SetFooter("Powered by ticketsbot.net", shard.SelfAvatar(256))
		}

		embed.SetTitle(panel.Title)
		embed.SetDescription(panel.Content)
		embed.SetColor(int(panel.Colour))

		msg, err := shard.CreateMessageEmbed(panel.ChannelId, embed); if err != nil {
			sentry.LogWithContext(err, errorContext)
			continue
		}

		// ReactionEmote is the unicode emoji
		if err = shard.CreateReaction(panel.ChannelId, msg.Id, panel.ReactionEmote); err != nil {
			sentry.LogWithContext(err, sentry.ErrorContext{})
		}

		//go dbclient.Client.Panel.Create(msg.Id, panel.ChannelId, panel.GuildId, panel.Title, panel.Content, panel.Colour, panel.TargetCategory, panel.ReactionEmote)
		go dbclient.Client.Panel.Create(database.Panel{
			MessageId:      msg.Id,
			ChannelId:      panel.ChannelId,
			GuildId:        panel.GuildId,
			Title:         	panel.Title,
			Content:        panel.Content,
			Colour:         panel.Colour,
			TargetCategory: panel.TargetCategory,
			ReactionEmote:  panel.ReactionEmote,
		})
	}
}
