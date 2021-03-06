package command

import (
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/database"
	"github.com/TicketsBot/TicketsGo/sentry"
	"github.com/rxdn/gdl/objects/channel/embed"
	"strings"
)

type RemoveSupportCommand struct {
}

func (RemoveSupportCommand) Name() string {
	return "removesupport"
}

func (RemoveSupportCommand) Description() string {
	return "Revokes a user's or role's support representative privileges"
}

func (RemoveSupportCommand) Aliases() []string {
	return []string{}
}

func (RemoveSupportCommand) PermissionLevel() utils.PermissionLevel {
	return utils.Admin
}

func (RemoveSupportCommand) Execute(ctx utils.CommandContext) {
	usageEmbed := embed.EmbedField{
		Name:   "Usage",
		Value:  "`t!removesupport @User`\n`t!removesupport @Role`\n`t!removesupport role name`",
		Inline: false,
	}

	if len(ctx.Args) == 0 {
		ctx.SendEmbed(utils.Red, "Error", "You need to mention a user or name a role to revoke support representative privileges from", usageEmbed)
		ctx.ReactWithCross()
		return
	}

	// get guild object
	guild, err := ctx.Guild(); if err != nil {
		sentry.ErrorWithContext(err, ctx.ToErrorContext())
		return
	}

	roles := make([]uint64, 0)
	if len(ctx.Message.Mentions) > 0 { // Individual users
		for _, mention := range ctx.Message.Mentions {
			// Verify that we're allowed to perform the remove operation
			if guild.OwnerId == mention.Id {
				ctx.SendEmbed(utils.Red, "Error", "The guild owner must be an admin")
				continue
			}

			if ctx.Author.Id == mention.Id {
				ctx.SendEmbed(utils.Red, "Error", "You cannot revoke your own privileges")
				continue
			}

			go func() {
				if err := database.Client.Permissions.RemoveSupport(ctx.GuildId, mention.Id); err != nil {
					sentry.ErrorWithContext(err, ctx.ToErrorContext())
					ctx.ReactWithCross()
				}
			}()
		}
	} else if len(ctx.Message.MentionRoles) > 0 {
		for _, mention := range ctx.Message.MentionRoles {
			roles = append(roles, mention)
		}
	} else { // Role
		roleName := strings.ToLower(strings.Join(ctx.Args, " "))

		// Get role ID from name
		valid := false
		for _, role := range guild.Roles {
			if strings.ToLower(role.Name) == roleName {
				roles = append(roles, role.Id)
				valid = true
				break
			}
		}

		// Verify a valid role was mentioned
		if !valid {
			ctx.SendEmbed(utils.Red, "Error", "You need to mention a user or name a role to revoke support representative privileges from", usageEmbed)
			ctx.ReactWithCross()
			return
		}
	}

	// Remove roles from DB
	for _, role := range roles {
		go func() {
			if err := database.Client.RolePermissions.RemoveSupport(ctx.GuildId, role); err != nil {
				sentry.ErrorWithContext(err, ctx.ToErrorContext())
				ctx.ReactWithCross()
			}
		}()
	}

	ctx.ReactWithCheck()
}

func (RemoveSupportCommand) Parent() interface{} {
	return nil
}

func (RemoveSupportCommand) Children() []Command {
	return make([]Command, 0)
}

func (RemoveSupportCommand) PremiumOnly() bool {
	return false
}

func (RemoveSupportCommand) Category() Category {
	return Settings
}

func (RemoveSupportCommand) AdminOnly() bool {
	return false
}

func (RemoveSupportCommand) HelperOnly() bool {
	return false
}
