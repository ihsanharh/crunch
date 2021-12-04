package music

import (
	"time"

	"github.com/XzFrosT/crunch/core"
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var VServerUpdateEvent = core.Event{
	Handler: func(e *gateway.VoiceServerUpdateEvent) {
		time.Sleep(1 * time.Second)

		if player := music.GetPlayer(e.GuildID); player != nil && player.State == music.PlayingState {
			player.Voicy.SendSpeaking()
		}
	},
}

var VStateUpdateEvent = core.Event{
	Handler: func(e *gateway.VoiceStateUpdateEvent) {
		if e.UserID != core.Self.ID {
			return
		}

		if player := music.GetPlayer(e.GuildID); player != nil && e.ChannelID.IsNull() {
			player.Stop(false)
		}
	},
}