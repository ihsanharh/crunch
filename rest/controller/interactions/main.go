package interactions

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/XzFrosT/crunch/core"
	"github.com/XzFrosT/crunch/utils/logger"
	"github.com/buger/jsonparser"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/gofiber/fiber/v2"
)

var decodedPubkey []byte

func Post(ctx *fiber.Ctx) error {
	timestamp := ctx.Request().Header.Peek("X-Signature-Timestamp")
	signature := ctx.Request().Header.Peek("X-Signature-Ed25519")

	if signature == nil || timestamp == nil {
		return fiber.ErrBadRequest
	}

	if !isValid(signature, append(timestamp, ctx.Body()...)) {
		return fiber.ErrUnauthorized
	}

	var event discord.InteractionEvent
	if err := ctx.BodyParser(&event); err != nil {
		if interactionType, _ := jsonparser.GetInt(ctx.Body(), "type"); interactionType == 1 {
			return ctx.JSON(api.InteractionResponse{Type: api.PongInteraction})
		}

		return err
	}

	if event.GuildID.IsValid() {
		event.User = &event.Member.User
	}

	if data, ok := event.Data.(*discord.CommandInteraction); ok {
		if command := core.Commands[data.Name]; command != nil {
			context := core.NewCommandContext(&event, core.State, data, ctx, command.Deffered)

			if command.Deffered {
				go command.Handler(context)
				return ctx.JSON(api.InteractionResponse{Type: api.DeferredMessageInteractionWithSource})
			}

			go command.Handler(context)
			return context.Wait()
		}
	}

	return fiber.ErrNotFound
}

func isValid(signature, hash []byte) bool {
	n, err := hex.Decode(signature, signature)
	if err != nil {
		return false
	}

	if decodedPubkey == nil {
		decodedPubkey, err = hex.DecodeString(core.App.VerifyKey)
		if err != nil {
			logger.Error("Unable to decode Discord's public key.", err)
			return false
		}
	}

	return ed25519.Verify(decodedPubkey, hash, signature[:n])
}
