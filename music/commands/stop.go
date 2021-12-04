package music

import (
	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/XzFrosT/crunch/utils/emojis"
)

var StopCommand = core.Command{
	Name:        "stop",
	Description: "Stop played music",
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.xmark, "Please join a voice channel first!")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil {
			ctx.Ephemeral().Reply(emojis.xmark, "There is nothing playing at the moment.")
			return
		}

		player.Stop(false)
		ctx.Reply(emojis.check, "Stopped music.")
	},
}
