package music

import (
	"github.com/XzFrosT/crunch/core"
	commands "github.com/XzFrosT/crunch/music/commands"
	events "github.com/XzFrosT/crunch/music/events"
)

var Module = &core.Module{
	Commands: []*core.Command{&commands.PlayCommand, &commands.SkipCommand, &commands.StopCommand, &commands.PauseCommand, &commands.ResumeCommand, &commands.SeekCommand, &commands.NowplayingCommand, &commands.ShuffleCommand},
	Events:   []*core.Event{&events.VServerUpdateEvent, &events.VStateUpdateEvent},
}