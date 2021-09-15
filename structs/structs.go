package structs

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

type Command struct {
	Name        string
	Description string
	Group       string
	Usage       string
	OwnerOnly   bool
	Run         func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData)
}
