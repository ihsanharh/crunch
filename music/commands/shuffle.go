package music

import (
	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/XzFrosT/crunch/utils/emojis"
)

var ShuffleCommand = core.Command{
	Name:        "shuffle",
	Description: "Shuffle the music in the queue",
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.xmark, "Please connect to a voice channel first!")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil || player.State == music.StoppedState {
			ctx.Ephemeral().Reply(emojis.xmark, "There is nothing playing at the moment.")
			return
		}

		if len(player.Queue) < 2 {
			ctx.Ephemeral().Reply(emojis.xmark, "Not enough songs to shuffle in the queue")
			return
		}

		player.Shuffle()
		ctx.Reply(emojis.check, "Successfully shuffled songs")
	},
}