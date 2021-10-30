package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Group       string
	Usage       string
	Options     []discord.CommandOption
	OwnerOnly   bool
	Run         func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData)
}

type Component struct {
	Run func(e *gateway.InteractionCreateEvent, data *discord.ComponentInteractionData)
}

var MapCommands = make(map[string]Command)
var MapComponents = make(map[string]Component)

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

	if len(commandFields) == 0 {
		commandFields = append(commandFields, discord.EmbedField{
			Name:  "There are no commands in this category!",
			Value: "that, or I wasn't able to grab them correctly :|\nEither way, if you have any command suggestions please enter them in \n<#689763311268397097>!",
		})
		return commandFields, utils.DiscordRed
	} else {
		return commandFields, utils.DiscordGreen
	}
}

func AddHandlers() {
	core.State.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch data := e.Data.(type) {
		case *discord.CommandInteractionData:
			if cmd, ok := MapCommands[data.Name]; ok {
				if cmd.OwnerOnly == true && e.Member.User.ID != discord.UserID(utils.MustSnowflakeEnv(core.Config.OwnerID)) {
					NoPerms(e, data, cmd)
					return
				}
				cmd.Run(e, data)
			}
		case *discord.ComponentInteractionData:
			if cmd, ok := MapComponents[data.CustomID]; ok {
				cmd.Run(e, data)
			}
		}
	})
}

func Register(appID discord.AppID, guildID discord.GuildID) {

	var commands []discord.Command

	for _, command := range MapCommands {
		commands = append(commands, discord.Command{
			Type:                discord.CommandType(1),
			Name:                command.Name,
			Description:         command.Description,
			Options:             command.Options,
			NoDefaultPermission: false,
		})
	}

	_, err := core.State.BulkOverwriteGuildCommands(appID, guildID, commands)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to overwrite commands in TDC with err: %v", err))
	}
}

func NoPerms(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData, cmd Command) {
	res := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{
				{
					Color:       utils.DiscordRed,
					Title:       "WOAH! You don't have permission to execute this command!",
					Description: fmt.Sprintf("Sorry, but %v has `Owneronly` set to %v.\n\nIf you believe this is an error please message <@148203660088705025>", cmd.Name, cmd.OwnerOnly),
				},
			},
		},
	}
	if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
		logger.Error(err)
	}
}
