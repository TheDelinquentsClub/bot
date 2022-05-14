package commands

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

func init() {
	MapModals["verification"] = Modal{
		Run: func(e *gateway.InteractionCreateEvent, data *discord.ModalInteraction) {

			arc := *data.Components[0].(*discord.ActionRowComponent)
			text := arc[0].(*discord.TextInputComponent).Value.Val

			var res api.InteractionResponse

			stmt, err := core.DB.Prepare("UPDATE players SET Discord_UUID = ?, status = 'LINKED' WHERE ID = ?")
			if err != nil {
				logger.Error("Failed to prepare query:", err)
			}

			_, err = stmt.Exec(e.Member.User.ID.String(), text)
			if err != nil {
				logger.Error("Failed to very account for", e.Member.User.ID.String(), "code:", text, "Error:", err)
				res = api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Color:       utils.DiscordRed,
								Title:       "Failed to verify this account!",
								Description: "Try again later, or contact <@148203660088705025> immediately in <#974302886945243156>",
								Timestamp:   discord.NowTimestamp(),
							},
						},
					},
				}
			} else {
				res = api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Flags: api.EphemeralResponse,
						Embeds: &[]discord.Embed{
							{
								Color: utils.DiscordGreen,
								Title: "Successfully connected your Minecraft account!",
								Description: "Thanks for connecting your account!\n" +
									"This will give you and other players access to certain commands that wouldn't work without this.\n" +
									"This also allows us to work on more commands for everyone!",
								Timestamp: discord.NowTimestamp(),
							},
						},
					},
				}
			}
			stmt.Close()
			if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
				logger.Error(err)
			}
		},
	}
}
