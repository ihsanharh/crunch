package core

import (
	"context"
	"reflect"

	"github.com/XzFrosT/crunch/utils/logger"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/pkg/errors"
)

type Event struct {
	Handler interface{}
}

type Module struct {
	Commands  []*Command
	Events    []*Event
	StartFunc func()
}

var (
	State *state.State
	Self  *discord.User
	App   *discord.Application

	Commands = make(map[string]*Command)
)

func NewClient(token string) (err error) {
	State, err = state.NewWithIntents("Bot "+token, gateway.IntentGuilds, gateway.IntentGuildVoiceStates)

	return err
}

func Connect() error {
	err := State.Open(context.Background())

	if err == nil {
		Self, err = State.Me()
	}
	if err == nil {
		App, err = State.CurrentApplication()
	}

	return err
}

func DeployCommands() error {

	previous, err := State.Commands(App.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to get Discord command list")
	}

	checked := make(map[string]interface{})
	for _, prevCmd := range previous {
		newCmd := Commands[prevCmd.Name]

		if newCmd == nil {
			logger.DebugF("Removed \"%s\" command from Discord.", prevCmd.Name)
			if err := State.DeleteCommand(App.ID, prevCmd.ID); err != nil {
				return errors.Wrapf(err, "failed to delete \"%s\" command", prevCmd.Name)
			}
		} else {
			if !reflect.DeepEqual(prevCmd.Options, newCmd.Options) || newCmd.Description != prevCmd.Description {
				logger.DebugF("Updating %s command in Discord.", newCmd.Name)
				if _, err := State.EditCommand(App.ID, prevCmd.ID, newCmd.RAW()); err != nil {
					return errors.Wrapf(err, "failed to update \"%s\" command", newCmd.Name)
				}
			}
			checked[newCmd.Name] = true
		}
	}

	for _, command := range Commands {
		if checked[command.Name] == nil {
			logger.DebugF("Creating %s command in Discord.", command.Name)
			if _, err := State.CreateCommand(App.ID, command.RAW()); err != nil {
				return errors.Wrapf(err, "failed to create \"%s\" command", command.Name)
			}
		}
	}

	return nil
}

func Close() {
	State.Close()
}

func AddModules(modules ...*Module) {
	for _, module := range modules {
		AddModule(module)
	}
}

func AddModule(module *Module) {
	for _, cmd := range module.Commands {
		cmd.Module = module
		Commands[cmd.Name] = cmd
	}

	for _, event := range module.Events {
		State.AddHandler(event.Handler)
	}

	if module.StartFunc != nil {
		go module.StartFunc()
	}
}
