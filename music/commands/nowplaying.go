package music

import (
	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/XzFrosT/crunch/utils"
	"github.com/XzFrosT/crunch/utils/emojis"
)

var NowplayingCommand = core.Command{
	Name:        "nowplaying",
	Description: "Shows what song is playing now",
	Handler: func(ctx *core.CommandContext) {
		player := music.GetPlayer(ctx.GuildID)

		if player == nil || player.State == music.StoppedState {
			ctx.Ephemeral().Reply(emojis.xmark, "There is nothing playing at the moment..")
			return
		}

		embed := utils.NewEmbed().
			Description("%s Now Playing: **[%s](%s)**", emojis.music, player.Current.Title, player.Current.URL).
			Thumbnail(player.Current.Thumbnail).
			Color(0x00FF59).
			Field("Author", player.Current.Author, true).
			Field("Durationo", utils.Fmt("%v/%v", utils.FormatTime(player.Voicy.Position), utils.Is(player.Current.IsLive, "--:--", utils.FormatTime(player.Current.Duration))), true).
			Field("Provider", player.Current.Provider(), true).
			Footer(utils.Fmt("Added by %s#%s", player.Current.Requester.Username, player.Current.Requester.Discriminator), player.Current.Requester.AvatarURL()).
			Timestamp(player.Current.RequestedAt)

		if player.State == music.PausedState {
			embed.Color(0xB4BE10).Description("%s Currently Paused on: [%s](%s)", emojis.xmark, player.Current.Title, player.Current.URL)
		}

		ctx.Reply(embed)
	},
}
