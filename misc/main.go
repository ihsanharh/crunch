package misc

import (
	"github.com/XzFrosT/crunch/core"
	commands "github.com/XzFrosT/crunch/misc/commands"
	events "github.com/XzFrosT/crunch/misc/events"
)

var Module = &core.Module{
	Commands: []*core.Command{commands.PingCommand},
	Events:   []*core.Event{events.ReadyEvent},
}
