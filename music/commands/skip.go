package music

import (
	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/XzFrosT/crunch/utils/emojis"
)

var SkipCommand = core.Command{
	Name:        "skip",
	Description: "Skip to the next song in queue",
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.xmark, "Please join a voice channel first!")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil || player.State == music.StoppedState {
			ctx.Ephemeral().Reply(emojis.xmark, "There is nothing playing at the moment")
			return
		}

		player.Skip()
		ctx.Reply(emojis.check, "Successfully skipped to the next song")
	},
}