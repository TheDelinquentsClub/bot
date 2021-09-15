package structs

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Run         func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData)
}
