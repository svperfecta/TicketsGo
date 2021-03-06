package command

import (
	"fmt"
	"github.com/TicketsBot/TicketsGo/bot/utils"
	"github.com/TicketsBot/TicketsGo/database"
	"github.com/TicketsBot/TicketsGo/sentry"
	"github.com/gofrs/uuid"
	"strconv"
	"strings"
	"time"
)

type AdminGeneratePremium struct {
}

func (AdminGeneratePremium) Name() string {
	return "genpremium"
}

func (AdminGeneratePremium) Description() string {
	return "Generate premium keys"
}

func (AdminGeneratePremium) Aliases() []string {
	return []string{"gp", "gk", "generatepremium", "genkeys", "generatekeys"}
}

func (AdminGeneratePremium) PermissionLevel() utils.PermissionLevel {
	return utils.Everyone
}

func (AdminGeneratePremium) Execute(ctx utils.CommandContext) {
	if len(ctx.Args) == 0 {
		ctx.ReactWithCross()
		return
	}

	days, err := strconv.Atoi(ctx.Args[0]); if err != nil {
		ctx.SendEmbed(utils.Red, "Admin", err.Error())
		ctx.ReactWithCross()
		return
	}

	amount := 1
	if len(ctx.Args) == 2 {
		if a, err := strconv.Atoi(ctx.Args[1]); err == nil {
			amount = a
		}
	}

	keys := make([]string, 0)
	for i := 0; i < amount; i++ {
		key, err := uuid.NewV4()
		if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
			continue
		}

		err = database.Client.PremiumKeys.Create(key, time.Hour * 24 * time.Duration(days))
		if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
		} else {
			keys = append(keys, key.String())
		}
	}

	dmChannel, err := ctx.Shard.CreateDM(ctx.Author.Id); if err != nil {
		ctx.SendEmbed(utils.Red, "Admin", err.Error())
		ctx.ReactWithCross()
		return
	}

	content := "```"
	for _, key := range keys {
		content += fmt.Sprintf("%s\n", key)
	}
	content = strings.TrimSuffix(content, "\n")
	content += "```"

	_, err = ctx.Shard.CreateMessage(dmChannel.Id, content); if err != nil {
		ctx.SendEmbed(utils.Red, "Admin", err.Error())
		ctx.ReactWithCross()
		return
	}

	ctx.ReactWithCheck()
}

func (AdminGeneratePremium) Parent() interface{} {
	return &AdminCommand{}
}

func (AdminGeneratePremium) Children() []Command {
	return []Command{}
}

func (AdminGeneratePremium) PremiumOnly() bool {
	return false
}

func (AdminGeneratePremium) Category() Category {
	return Settings
}

func (AdminGeneratePremium) AdminOnly() bool {
	return true
}

func (AdminGeneratePremium) HelperOnly() bool {
	return false
}
