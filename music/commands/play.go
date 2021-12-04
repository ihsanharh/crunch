package music

import (
	"github.com/XzFrosT/crunch/core"

	music "github.com/XzFrosT/crunch/music/audio"
	providers "github.com/XzFrosT/crunch/music/providers"
	"github.com/XzFrosT/crunch/utils"
	"github.com/XzFrosT/crunch/utils/emojis"
	"github.com/diamondburned/arikawa/v3/discord"
)

var PlayCommand = core.Command{
	Name:        "play",
	Description: "Play a song",
	Options: discord.CommandOptions{&discord.StringOption{
		OptionName:  "song",
		Description: "Name of the song/playlist or the URL(only youtube)",
		Required:    true,
	}, &discord.BooleanOption{
		OptionName:  "shuffle",
		Description: "Shuffle the music in the queue",
	}},
	Handler: func(ctx *core.CommandContext) {
		query, shuffle := ctx.Argument(0).String(), ctx.Argument(1).Bool()

		state := ctx.VoiceState()
		if state == nil {
			ctx.Ephemeral().Reply(emojis.xmark, "Please connect to a voice channel first!")
			return
		}

		embed := utils.NewEmbed().Color(0xF0FF00).Description("%s Getting results for your search...", emojis.search)
		go ctx.Reply(embed)

		player := music.GetOrCreatePlayer(ctx.GuildID, ctx.ChannelID, state.ChannelID)
		defer checkIdle(player)

		result, err := providers.FindSong(query)
		if err != nil {
			ctx.Stacktrace(err)
			return
		}

		if result == nil {
			ctx.Reply(embed.Color(0xF93A2F).Description("%s I couldn't find this song.", emojis.xmark))
			return
		}

		if result.Playlist != nil {
			player.AddSong(ctx.User, shuffle, result.Songs...)

			ctx.Reply(embed.Color(0x00D166).
				Description("%s Playlist [%s](%s) added to queue", emojis.check, result.Playlist.Title, result.Playlist.URL).
				Field("Creator", result.Playlist.Author, true).
				Field("Songs", len(result.Songs), true).
				Field("Duration", utils.FormatTime(result.Playlist.Duration), true))
			return
		}

		song := result.Songs[0]
		embed.Thumbnail(song.Thumbnail).
			Field("Author", song.Author, true).
			Field("Duration", utils.Is(song.IsLive, "--:--", utils.FormatTime(song.Duration)), true).
			Field("Provider", song.Provider(), true)

		if !song.IsLoaded() {
			go ctx.Reply(embed.Description("%s Loading [%s](%s)", emojis.AnimatedStaff, song.Title, song.URL))

			if err := song.Load(); err != nil {
				ctx.Stacktrace(err)
				return
			}
		}

		player.AddSong(ctx.User, shuffle, song)
		ctx.Reply(embed.
			Color(0x00D166).
			Thumbnail(song.Thumbnail).
			Description("%s Music [%s](%s) added to queue", emojis.Yeah, song.Title, song.URL))
	},
}

func checkIdle(player *music.Player) {
	if player.State != music.StoppedState || len(player.Queue) != 0 {
		return
	}

	player.Stop(true)
}
