package bot

import (
	"github.com/bwmarrin/discordgo"
)

func NewTodoEventEmbed(topic, event string) *discordgo.MessageEmbed {
	return NewEmbed("Todo Event", "", 0x00ff00,
		&discordgo.MessageEmbedField{
			Name:   topic,
			Value:  event,
			Inline: false,
		})
}

func NewErrorEmbed(errMsg string) *discordgo.MessageEmbed {
	return NewEmbed("Error", "", 0xff0000,
		&discordgo.MessageEmbedField{
			Name:   "error",
			Value:  errMsg,
			Inline: false,
		})
}

func NewEmbed(title, description string, color int, fields ...*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Fields:      fields,
	}
}
