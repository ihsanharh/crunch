package core

import (
	"github.com/XzFrosT/crunch/utils"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type BasicContext struct {
	ChannelID discord.ChannelID
	GuildID   discord.GuildID
	State     *state.State

	response *api.SendMessageData
}

func NewBasicContext(channelID discord.ChannelID, guildID discord.GuildID) *BasicContext {
	return &BasicContext{ChannelID: channelID, GuildID: guildID, State: State, response: &api.SendMessageData{}}
}

func (ctx *BasicContext) Guild() (*discord.Guild, error) {
	if ctx.GuildID.IsNull() {
		return nil, nil
	}

	return ctx.State.Guild(ctx.GuildID)
}

func (ctx *BasicContext) File(file sendpart.File) {
	ctx.response.Files = append(ctx.response.Files, file)
}

func (ctx *BasicContext) Embed(embed *utils.Embed) {
	ctx.response.Embeds = []discord.Embed{embed.Build()}
}

func (ctx *BasicContext) Send(args ...interface{}) {
	if len(args) > 1 {
		ctx.response.Content = utils.Fmt("%v | %v", args[0], utils.Fmt(args[1].(string), args[2:]...))
	}

	if len(args) == 1 {
		if embed, ok := args[0].(*utils.Embed); ok {
			ctx.Embed(embed)
		} else {
			ctx.response.Content = utils.Fmt("%v", args[0])
		}
	}

	ctx.State.SendMessageComplex(ctx.ChannelID, *ctx.response)
	ctx.response = &api.SendMessageData{}
}
