package command

import (
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/database"
	"github.com/TicketsBot/TicketsGo/sentry"
	"strconv"
	"strings"
)

type ManageCannedResponsesAdd struct {
}

func (ManageCannedResponsesAdd) Name() string {
	return "add"
}

func (ManageCannedResponsesAdd) Description() string {
	return "Adds a new canned response"
}

func (ManageCannedResponsesAdd) Aliases() []string {
	return []string{"new"}
}

func (ManageCannedResponsesAdd) PermissionLevel() utils.PermissionLevel {
	return utils.Support
}

func (ManageCannedResponsesAdd) Execute(ctx utils.CommandContext) {
	guildId, err := strconv.ParseInt(ctx.Guild.ID, 10, 64); if err != nil {
		sentry.ErrorWithContext(err, ctx.ToErrorContext())
		return
	}

	if len(ctx.Args) < 2 {
		ctx.ReactWithCross()
		ctx.SendEmbed(utils.Red, "Error", "You must specify a canned response ID and the contents of the response")
		return
	}

	id := ctx.Args[0]
	content := ctx.Args[1:] // content cannot be bigger than the Discord limit, obviously

	if len(id) > 16 {
		ctx.ReactWithCross()
		ctx.SendEmbed(utils.Red, "Error", "A canned response with the ID `$id` already exists. You can delete the response using `t!mcr delete $id`")
		return
	}

	go database.AddCannedResponse(guildId, id, strings.Join(content, " "))
	ctx.ReactWithCheck()
}

func (ManageCannedResponsesAdd) Parent() interface{} {
	return ManageCannedResponses{}
}

func (ManageCannedResponsesAdd) Children() []Command {
	return make([]Command, 0)
}

func (ManageCannedResponsesAdd) PremiumOnly() bool {
	return false
}

func (ManageCannedResponsesAdd) AdminOnly() bool {
	return false
}

func (ManageCannedResponsesAdd) HelperOnly() bool {
	return false
}
