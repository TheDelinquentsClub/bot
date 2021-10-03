package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/structs"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
)

var MapCommands = make(map[string]structs.Command)
var MapComponents = make(map[string]structs.Component)

// GetCommands returns all available commands that have been mapped for use in /help
func GetCommands(category string) ([]discord.EmbedField, discord.Color) {

	var commandFields []discord.EmbedField

	for _, command := range MapCommands {
		switch command.Group {
		case category:
			commandFields = append(commandFields, discord.EmbedField{
				Name:   strings.Title(command.Name),
				Value:  fmt.Sprintf("%v\nUsage: `%v`\nOwner Only: %v", command.Description, command.Usage, command.OwnerOnly),
				Inline: true,
			})
		default:
			continue
		}
	}

	if fmt.Sprint(commandFields) == "[]" {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "There are no commands in this category!",
			Value: "that, or I wasn't able to grab them correctly :|\nEither way, if you have any command suggestions please enter them in \n<#689763311268397097>!",
		})
		return commandFields, utils.DiscordRed
	} else {
		return commandFields, utils.DiscordGreen
	}
}
