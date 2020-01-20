package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/TicketsGo/config"
	"github.com/TicketsBot/TicketsGo/database"
	"github.com/TicketsBot/TicketsGo/sentry"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ProxyResponse struct {
	Premium bool
	Tier int
}

var premiumCache = cache.New(10 * time.Minute, 10 * time.Minute)

func IsPremiumGuild(ctx CommandContext, ch chan bool) {
	premium, ok := premiumCache.Get(ctx.Guild)

	if ok {
		ch<-premium.(bool)
		return
	}

	// First lookup by premium key, then votes, then patreon
	keyLookup := make(chan bool)
	go database.IsPremium(ctx.GuildId, keyLookup)

	if <-keyLookup {
		err := premiumCache.Add(ctx.Guild, true, 10 * time.Minute)

		if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
		}

		ch<-true
	} else {
		// Get guild object
		guild, err := ctx.Session.State.Guild(ctx.Guild); if err != nil {
			guild, err = ctx.Session.Guild(ctx.Guild); if err != nil {
				sentry.ErrorWithContext(err, ctx.ToErrorContext())
				ch<-false
				return
			}
		}

		// Lookup votes
		ownerId, err := strconv.ParseInt(guild.OwnerID, 10, 64); if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
			ch <- false
			return
		}

		hasVoted := make(chan bool)
		go database.HasVoted(ownerId, hasVoted)
		if <-hasVoted {
			ch <- true

			err = premiumCache.Add(ctx.Guild, true, 10 * time.Minute)

			if err != nil {
				sentry.ErrorWithContext(err, ctx.ToErrorContext())
			}

			return
		}

		// Lookup Patreon
		client := &http.Client{
			Timeout: time.Second * 3,
		}

		url := fmt.Sprintf("%s/ispremium?key=%s&id=%s", config.Conf.Bot.PremiumLookupProxyUrl, config.Conf.Bot.PremiumLookupProxyKey, guild.OwnerID)
		req, err := http.NewRequest("GET", url, nil)

		res, err := client.Do(req); if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
			ch<-false
			return
		}
		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body); if err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
			ch<-false
			return
		}

		var proxyResponse ProxyResponse
		if err = json.Unmarshal(content, &proxyResponse); err != nil {
			sentry.ErrorWithContext(err, ctx.ToErrorContext())
			ch<-false
			return
		}

		// I think we can safely ignore this error as it's caused by a race condition
		// which doesn't have any negative effects - it'd mean we have to lock the entire map
		// while performing a lookup
		_ = premiumCache.Add(ctx.Guild, proxyResponse.Premium, 10 * time.Minute)

		ch <-proxyResponse.Premium
	}
}
