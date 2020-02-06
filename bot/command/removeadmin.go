package command

import (
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/database"
	"strings"
)

type RemoveAdminCommand struct {
}

func (RemoveAdminCommand) Name() string {
	return "removeadmin"
}

func (RemoveAdminCommand) Description() string {
	return "Revokes a user's or role's admin privileges"
}

func (RemoveAdminCommand) Aliases() []string {
	return []string{}
}

func (RemoveAdminCommand) PermissionLevel() utils.PermissionLevel {
	return utils.Admin
}

func (RemoveAdminCommand) Execute(ctx utils.CommandContext) {
	if len(ctx.Args) == 0 {
		ctx.SendEmbed(utils.Red, "Error", "You need to mention a user or name a role to revoke admin privileges from")
		ctx.ReactWithCross()
		return
	}

	var roleId string
	if len(ctx.Message.Mentions) > 0 {
		for _, mention := range ctx.Message.Mentions {
			if ctx.Guild.OwnerID == mention.ID {
				ctx.SendEmbed(utils.Red, "Error", "The guild owner must be an admin")
				continue
			}

			if ctx.User.ID == mention.ID {
				ctx.SendEmbed(utils.Red, "Error", "You cannot revoke your own privileges")
				continue
			}

			go database.RemoveAdmin(ctx.Guild.ID, mention.ID)
		}
	} else {
		roleName := strings.ToLower(ctx.Args[0])

		// Get role ID from name
		for _, role := range ctx.Guild.Roles {
			if strings.ToLower(role.Name) == roleName {
				roleId = role.ID
				break
			}
		}

		// Verify a valid role was mentioned
		if roleId == "" {
			ctx.SendEmbed(utils.Red, "Error", "You need to mention a user or name a role to revoke admin privileges from")
			ctx.ReactWithCross()
			return
		}

		go database.RemoveAdminRole(ctx.Guild.ID, roleId)
	}

	ctx.ReactWithCheck()
}

func (RemoveAdminCommand) Parent() interface{} {
	return nil
}

func (RemoveAdminCommand) Children() []Command {
	return make([]Command, 0)
}

func (RemoveAdminCommand) PremiumOnly() bool {
	return false
}

func (RemoveAdminCommand) AdminOnly() bool {
	return false
}

func (RemoveAdminCommand) HelperOnly() bool {
	return false
}
