package utils

import (
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

type Embed struct {
	discord.Embed
}

func NewEmbed() *Embed {
	return &Embed{discord.Embed{}}
}

func (e *Embed) Author(args ...string) *Embed {
	author := &discord.EmbedAuthor{Name: args[0]}
	if len(args) >= 2 {
		author.Icon = args[1]
	}
	if len(args) >= 3 {
		author.URL = args[2]
	}

	e.Embed.Author = author
	return e
}

func (e *Embed) URL(url string) *Embed {
	e.Embed.URL = url
	return e
}

func (e *Embed) Title(title string, args ...interface{}) *Embed {
	e.Embed.Title = Fmt(title, args...)
	return e
}

func (e *Embed) Description(desc string, args ...interface{}) *Embed {
	e.Embed.Description = Fmt(desc, args...)
	return e
}

func (e *Embed) Color(color int) *Embed {
	e.Embed.Color = discord.Color(color)
	return e
}

func (e *Embed) Image(url string) *Embed {
	e.Embed.Image = &discord.EmbedImage{
		URL: url,
	}

	return e
}

func (e *Embed) Thumbnail(url string) *Embed {
	e.Embed.Thumbnail = &discord.EmbedThumbnail{
		URL: url,
	}

	return e
}

func (e *Embed) Timestamp(time time.Time) *Embed {
	e.Embed.Timestamp = discord.NewTimestamp(time)
	return e
}

func (e *Embed) Field(name, value interface{}, inline bool) *Embed {
	e.Fields = append(e.Fields, discord.EmbedField{
		Name:   Fmt("%v", name),
		Value:  Fmt("%v", value),
		Inline: inline,
	})
	return e
}

func (e *Embed) SetField(index int, name, value string, inline bool) *Embed {
	if index > len(e.Fields) {
		return e
	}

	e.Fields[index] = discord.EmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	return e
}

func (e *Embed) Footer(content, imgUrl string) *Embed {
	e.Embed.Footer = &discord.EmbedFooter{
		Text: content,
		Icon: imgUrl,
	}
	return e
}

func (e *Embed) Build() discord.Embed {
	return e.Embed
}