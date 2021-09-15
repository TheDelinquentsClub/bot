package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/structs"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
)

var CommandsMap = make(map[string]structs.Command)

var commands []discord.Command

func init() {
	for _, command := range CommandsMap {
		var id int64 = 0
		for i := 0; i < len(CommandsMap); i++ {
			id++
		}
		commands = append(commands, discord.Command{
			Name:                command.Name,
			Description:         command.Description,
			NoDefaultPermission: command.OwnerOnly,
		})
		if command.OwnerOnly == true {
			core.State.BatchEditCommandPermissions(discord.AppID(utils.MustSnowflakeEnv(core.Config.APPID)), discord.GuildID(utils.MustSnowflakeEnv(core.Config.GUILDID)), []api.BatchEditCommandPermissionsData{
				{
					ID: discord.CommandID(discord.id)),
					Permissions: []discord.CommandPermissions{
				{
					ID:utils.MustSnowflakeEnv(core.Config.OwnerID),
					Type: 2,
					Permission: true,
				},
				},
				},
			})
		}
	}
}

// GetCommands returns all available commands that have been mapped for use in /help
func GetCommands() []discord.EmbedField {
	var infoGroupString []string
	var funGroupString []string
	var debugGroupString []string

	for _, command := range CommandsMap {
		switch command.Group {
		case "info":
			infoGroupString = append(infoGroupString, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		case "fun":
			funGroupString = append(funGroupString, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		case "debug":
			debugGroupString = append(debugGroupString, fmt.Sprintf("`%s` - %s\n%s", command.Name, command.Description, command.Usage))
		default:
			continue // or have a default group
		}
	}

	var commandFields []discord.EmbedField

	if infoGroupString != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Info",
			Value: strings.Join(infoGroupString, "\n"),
		})
	}
	if funGroupString != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Fun",
			Value: strings.Join(funGroupString, "\n"),
		})
	}
	if debugGroupString != nil {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "Debug",
			Value: strings.Join(debugGroupString, "\n"),
		})
	}

	return commandFields
}
