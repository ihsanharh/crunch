package core

import (
	"github.com/XzFrosT/crunch/utils"
	"github.com/XzFrosT/crunch/utils/emojis"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/gofiber/fiber/v2"
)

type Command struct {
	Name, Description string
	Module            *Module
	Deffered          bool
	Type              discord.CommandType
	Options           discord.CommandOptions
	Handler           func(*CommandContext)
}

type CommandContext struct {
	*discord.InteractionEvent

	Data  *discord.CommandInteraction
	State *state.State

	response *api.InteractionResponseData
	sended   bool

	request     *fiber.Ctx
	requestChan chan error
}

type Argument struct {
	a *discord.CommandInteractionOption
}

func NewCommandContext(e *discord.InteractionEvent, state *state.State, data *discord.CommandInteraction, requestCtx *fiber.Ctx, sended bool) *CommandContext {
	return &CommandContext{
		InteractionEvent: e,
		State:            state,
		Data:             data,
		response:         &api.InteractionResponseData{},
		request:          requestCtx,
		requestChan:      make(chan error),
		sended:           sended,
	}
}

func (ctx *CommandContext) Wait() error {
	defer close(ctx.requestChan)
	return <-ctx.requestChan
}

func (ctx *CommandContext) Argument(index int) *Argument {
	if len(ctx.Data.Options) < index+1 {
		return &Argument{}
	}

	return &Argument{a: &ctx.Data.Options[index]}
}

func (ctx *CommandContext) Guild() *discord.Guild {
	if !ctx.GuildID.IsValid() {
		return nil
	}

	guild, _ := ctx.State.Guild(ctx.GuildID)
	return guild
}

func (ctx *CommandContext) VoiceState() *discord.VoiceState {
	if !ctx.GuildID.IsValid() {
		return nil
	}

	state, _ := ctx.State.VoiceState(ctx.GuildID, ctx.Member.User.ID)
	return state
}

func (ctx *CommandContext) File(file sendpart.File) *CommandContext {
	ctx.response.Files = append(ctx.response.Files, file)
	return ctx
}

func (ctx *CommandContext) Embed(eb *utils.Embed) *CommandContext {
	ctx.response.Embeds = &[]discord.Embed{eb.Embed}
	return ctx
}

func (ctx *CommandContext) Ephemeral() *CommandContext {
	ctx.response.Flags = 1 << 6
	return ctx
}

func (ctx *CommandContext) Reply(args ...interface{}) {
	ctx.checkArguments(args...)

	if ctx.sended {
		ctx.Edit(args...)
		return
	}

	ctx.sended = true
	ctx.requestChan <- ctx.request.JSON(api.InteractionResponse{Type: api.MessageInteractionWithSource, Data: ctx.response})
}

func (ctx *CommandContext) Edit(args ...interface{}) (*discord.Message, error) {
	ctx.checkArguments(args...)

	return ctx.State.EditInteractionResponse(ctx.AppID, ctx.Token, api.EditInteractionResponseData{
		Content: ctx.response.Content, Components: ctx.response.Components,
		Embeds: ctx.response.Embeds, Files: ctx.response.Files,
	})
}

func (ctx *CommandContext) checkArguments(args ...interface{}) {
	if len(args) > 1 {
		ctx.response.Content = option.NewNullableString(utils.Fmt("%v | %v", args[0], utils.Fmt(args[1].(string), args[2:]...)))
	}

	if len(args) == 1 {
		if embed, ok := args[0].(*utils.Embed); ok {
			ctx.Embed(embed)
		} else {
			ctx.response.Content = option.NewNullableString(utils.Fmt("%v", args[0]))
		}
	}
}

func (ctx *CommandContext) Stacktrace(err error) {
	ctx.Reply(emojis.xmark, "An error occurred while performing this action.: ```go\n%+v```", err)
}

func (cmd *Command) RAW() api.CreateCommandData {
	return api.CreateCommandData{Name: cmd.Name, Description: cmd.Description, Type: cmd.Type, Options: cmd.Options}
}

func (argument *Argument) Bool() bool {
	if argument.a == nil {
		return false
	}

	value, _ := argument.a.BoolValue()
	return value
}

func (argument *Argument) String() string {
	if argument.a == nil {
		return ""
	}
	return argument.a.String()
}
