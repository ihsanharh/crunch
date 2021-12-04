package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/XzFrosT/crunch/core"
	"github.com/XzFrosT/crunch/misc"
	"github.com/XzFrosT/crunch/music"
	"github.com/XzFrosT/crunch/rest"
	"github.com/XzFrosT/crunch/utils/logger"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.Fatal("Failed to load environment variable", err)
		}
	}

	if err := core.NewClient(os.Getenv("TOKEN")); err != nil {
		logger.Fatal("Failed to build", err)
	}

	core.AddModules(music.Module, misc.Module, rest.Module)

	if err := core.Connect(); err != nil {
		logger.Fatal("Something went wrong when trying to connect to discord")
	}

	if err := core.DeployCommands(); err != nil {
		logger.Fatal("Something went wrong when trying to register command.", err)
	}

	logger.Info("Bot is now ready! ude CTRL + C to stop the bot.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
	core.Close()
}