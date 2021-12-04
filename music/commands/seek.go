package music

import (
	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/XzFrosT/crunch/utils"
	"github.com/XzFrosT/crunch/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var SeekCommand = core.Command{
	Name:        "seek",
	Description: "Seek current song to specific timestamp",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "position",
		Description: "Desired position, example of valid formats: 05:05 or 5m5s",
		Required:    true,
	}},
	Handler: func(ctx *core.CommandContext) {
		if ctx.VoiceState() == nil {
			ctx.Ephemeral().Reply(emojis.xmark, "Please connect to a voice channel first!")
			return
		}

		player := music.GetPlayer(ctx.GuildID)
		if player == nil || player.State != music.PlayingState {
			ctx.Ephemeral().Reply(emojis.xmark, "There is nothing playing at the moment or the music is paused.")
			return
		}

		if player.Current.IsLive {
			ctx.Ephemeral().Reply(emojis.xmark, "You cannot do this on live streams.")
			return
		}

		duration, err := utils.ParseDuration(ctx.Argument(0).String())
		if err != nil || duration < 0 || duration > player.Current.Duration {
			ctx.Ephemeral().Reply(emojis.xmark, "Invalid duration or greater than the total duration of the song.")
			return
		}

		player.Voicy.Seek(duration)
		ctx.Ephemeral().Reply(emojis.check, "Song timestamp changed to minutes `%s`.", utils.FormatTime(duration))
	},
}
