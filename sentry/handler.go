package sentry

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"runtime/debug"
)

type ErrorContext struct {
	Guild       uint64
	User        uint64
	Channel     uint64
	Shard       int
	Command     string
	Premium     bool
}

func Error(e error) {
	fmt.Println(e.Error())
	debug.PrintStack()
	/*wrapped := errors.New(e)
	raven.Capture(ConstructErrorPacket(wrapped), nil)*/
}

func LogWithContext(e error, ctx ErrorContext) {
	fmt.Println(e.Error())
	debug.PrintStack()
	/*wrapped := errors.New(e)
	raven.Capture(ConstructPacket(wrapped, raven.INFO), map[string]string{
		"guild":       strconv.FormatUint(ctx.Guild, 10),
		"user":        strconv.FormatUint(ctx.User, 10),
		"channel":     strconv.FormatUint(ctx.Channel, 10),
		"shard":       strconv.Itoa(ctx.Shard),
		"command":     ctx.Command,
		"premium":     strconv.FormatBool(ctx.Premium),
	})*/
}

func LogRestRequest(url string) {
	raven.CaptureMessage(url, nil, nil)
}

func ErrorWithContext(e error, ctx ErrorContext) {
	fmt.Println(e.Error())
	debug.PrintStack()
	/*wrapped := errors.New(e)
	raven.Capture(ConstructErrorPacket(wrapped), map[string]string{
		"guild":       strconv.FormatUint(ctx.Guild, 10),
		"user":        strconv.FormatUint(ctx.User, 10),
		"channel":     strconv.FormatUint(ctx.Channel, 10),
		"shard":       strconv.Itoa(ctx.Shard),
		"command":     ctx.Command,
		"premium":     strconv.FormatBool(ctx.Premium),
	})*/
}
