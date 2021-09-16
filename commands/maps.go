package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/structs"
	"strings"
)

var CommandsMap = make(map[string]structs.Command)

var commands []discord.Command

func init() {
	for _, command := range CommandsMap {
		commands = append(commands, discord.Command{
			Name:                command.Name,
			Description:         command.Description,
			Options:             command.Options,
			NoDefaultPermission: command.OwnerOnly,
		})
	}
}

// GetCommands returns all available commands that have been mapped for use in /help
func GetCommands() []discord.EmbedField {
	var infoGroup []string
	var funGroup []string
	var debugGroup []string
	var utilityGroup []string

	for _, command := range CommandsMap {
		switch command.Group {
		case "info":
			infoGroup = append(infoGroup, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		case "fun":
			funGroup = append(funGroup, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		case "debug":
			debugGroup = append(debugGroup, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		case "utility":
			utilityGroup = append(utilityGroup, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		default:
			continue // or have a default group
		}
	}

	var commandFields []discord.EmbedField

	if infoGroup != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Info",
			Value: strings.Join(infoGroup, "\n"),
		})
	}
	if funGroup != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Fun",
			Value: strings.Join(funGroup, "\n"),
		})
	}
	if debugGroup != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Debug",
			Value: strings.Join(debugGroup, "\n"),
		})
	}
	if utilityGroup != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Utility",
			Value: strings.Join(utilityGroup, "\n"),
		})
	}

	return commandFields
}
