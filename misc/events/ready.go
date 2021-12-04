package misc

import (
	"os"

	"github.com/XzFrosT/crunch/core"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var ReadyEvent = &core.Event{
	Handler: func(_ *gateway.ReadyEvent) {
		core.State.UpdateStatus(gateway.UpdateStatusData{Activities: []discord.Activity{{Name: os.Getenv("DISCORD_STATUS"), Type: discord.ListeningActivity}}})
	},
}
