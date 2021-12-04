package music 

import (
	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/XzFrosT/crunch/utils/emojis"
)

var PauseCommand = core.Command{
	Name:        "pause",
	Description: "Pause currently playing song",
	Handler:     func(ctx *core.CommandContext) { handleCommand(ctx) },
}

var ResumeCommand = core.Command{
	Name:        "resume",
	Description: "Resume paused queue",
	Handler:     func(ctx *core.CommandContext) { handleCommand(ctx) },
}

func handleCommand(ctx *core.CommandContext) {
	if ctx.VoiceState() == nil {
		ctx.Ephemeral().Reply(emojis.xmark, "Please connect to a voice channel first!")
		return
	}

	player := music.GetPlayer(ctx.GuildID)
	if player == nil || player.State == music.StoppedState {
		ctx.Ephemeral().Reply(emojis.xmark, "There is nothing playing at the moment.")
		return
	}

	if player.Current.IsLive {
		ctx.Ephemeral().Reply(emojis.xmark, "You cannot do this on live streams.")
		return
	}

	if player.State == music.PlayingState {
		player.Pause()
		ctx.Reply(emojis.check, "Music paused.")
	} else {
		player.Resume()
		ctx.Reply(emojis.check, "Resuming music...")
	}
}
