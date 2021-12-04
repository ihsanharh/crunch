package player

import (
	music "github.com/XzFrosT/crunch/music/audio"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	guildId, _ := discord.ParseSnowflake(c.Params("id"))
	player := music.GetPlayer(discord.GuildID(guildId))

	if !guildId.IsValid() || player == nil || player.Voicy == nil {
		return fiber.ErrNotFound
	}

	return c.JSON(&fiber.Map{"data": &fiber.Map{
		"current":  player.Current,
		"queue":    player.Queue,
		"state":    player.State,
		"position": player.Voicy.Position,
	}, "error": nil})
}
