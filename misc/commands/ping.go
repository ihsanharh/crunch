package misc

import (
	"time"

	"github.com/XzFrosT/crunch/core"
	"github.com/XzFrosT/crunch/utils/emojis"
)

var PingCommand = &core.Command{
	Name:        "ping",
	Description: "Get bot latency",
	Handler: func(ctx *core.CommandContext) {
		latency := time.Duration(ctx.State.PacerLoop.EchoBeat.Get() - ctx.State.PacerLoop.SentBeat.Get())
		if latency <= 0 {
			ctx.Reply(emojis.PingPong, "There aren't enough latency measurements yet ;(")
			return
		}

		ctx.Reply(emojis.PingPong, "Pong, %dms.", latency.Milliseconds())
	},
}
