package commands

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

// Command defines the basic structure of commands.
type Command struct {
	Type        discord.CommandType
	Name        string
	Description string
	Group       string
	Usage       string
	Options     []discord.CommandOption
	Restricted  bool
	Run         func(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction)
}

type Button struct {
	Run func(e *gateway.InteractionCreateEvent, data *discord.ButtonInteraction)
}
type Select struct {
	Run func(e *gateway.InteractionCreateEvent, data *discord.SelectInteraction)
}

var MapCommands = make(map[string]Command)
var MapButtons = make(map[string]Button)
var MapSelect = make(map[string]Select)

// GetCommands returns all available commands that have been mapped for use in /help
func GetCommands(category string) ([]discord.EmbedField, discord.Color) {

	var commandFields []discord.EmbedField

	for _, command := range MapCommands {
		switch command.Group {
		case category:
			commandFields = append(commandFields, discord.EmbedField{
				Name:   strings.Title(command.Name),
				Value:  fmt.Sprintf("%v\nUsage: `%v`\nOwner Only: %v", command.Description, command.Usage, command.Restricted),
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
		return commandFields, utils.DefaultColour
	}
}

func AddHandlers() {
	core.State.AddHandler(func(e *gateway.InteractionCreateEvent) {
		switch data := e.Data.(type) {
		case *discord.CommandInteraction:
			if cmd, ok := MapCommands[data.Name]; ok {
				if cmd.Restricted == true && e.Member.User.ID != discord.UserID(utils.MustSnowflakeEnv(core.Config.CreatorID)) {
					NoPerms(e, data, cmd)
					return
				}
				cmd.Run(e, data)
			} else {
				doesntExist(e, data)
			}
		case *discord.SelectInteraction:
			if cmd, ok := MapSelect[string(data.CustomID)]; ok {
				cmd.Run(e, data)
			}
		case *discord.ButtonInteraction:
			if cmd, ok := MapButtons[string(data.CustomID)]; ok {
				cmd.Run(e, data)
			}
		}
	})
}

func Register(appID discord.AppID, guildID discord.GuildID) {

	var commands []api.CreateCommandData

	for _, command := range MapCommands {
		if command.Type == 0 {
			command.Type = 1
		}
		commands = append(commands, api.CreateCommandData{
			Type:                command.Type,
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

func NoPerms(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction, cmd Command) {
	res := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Flags: api.EphemeralResponse,
			Embeds: &[]discord.Embed{
				{
					Color:       utils.DiscordRed,
					Title:       "WOAH! You don't have permission to execute this command!",
					Description: fmt.Sprintf("Sorry, but %v has `Owneronly` set to %v.\n\nIf you believe this is an error please message <@148203660088705025>", cmd.Name, cmd.Restricted),
				},
			},
		},
	}
	if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
		logger.Error(err)
	}
}
func doesntExist(e *gateway.InteractionCreateEvent, data *discord.CommandInteraction) {
	res := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Flags: api.EphemeralResponse,
			Embeds: &[]discord.Embed{
				{
					Color:       utils.DiscordRed,
					Title:       "wait a minute... this command doesnt exist!",
					Description: fmt.Sprintf("sorry but %v is not a valid command.\n\nIf you believe this is an error please message <@148203660088705025>", e.Message.Content),
				},
			},
		},
	}
	if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
		logger.Error(err)
	}
}
