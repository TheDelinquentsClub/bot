package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/structs"
	"github.com/kingultron99/tdcbot/utils"
)

func init() {
	MapCommands["ban"] = structs.Command{
		Name:        "ban",
		Description: "Bans a user",
		Usage:       "/ban <user>",
		Group:       "utility",
		Options: []discord.CommandOption{
			{
				Name:        "user",
				Description: "The user to ban",
				Type:        discord.CommandOptionType(6),
				Required:    true,
			},
			{
				Name:        "reason",
				Description: "Why is the user being banned?",
				Type:        discord.CommandOptionType(3),
				Required:    false,
			},
		},
		OwnerOnly: true,
		Run: func(e *gateway.InteractionCreateEvent, data *discord.CommandInteractionData) {
			var snowflake, _ = data.Options[0].Snowflake()

			if snowflake != utils.MustSnowflakeEnv(fmt.Sprint(e.Member.User.ID)) {
				var user, err = core.State.User(discord.UserID(snowflake))
				if err != nil {
					return
				}
				var reason string

				if len(data.Options) > 1 {
					reason = fmt.Sprint(data.Options[1])
				} else {
					reason = "No reason specified"
				}

				var banData = api.BanData{AuditLogReason: api.AuditLogReason(reason)}

				err = core.State.Ban(e.GuildID, user.ID, banData)
				if err != nil {
					logger.Error.Println(fmt.Sprintf("Failed to ban user %v\nerror: %v", user.Username, err))
				}

				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Title: "🔨 Banned!",
								Fields: []discord.EmbedField{
									{
										Name:  "User",
										Value: fmt.Sprintf("<@%v>", user.ID),
									},
									{
										Name:   "Reason",
										Value:  reason,
										Inline: true,
									},
								},
								Color:     utils.DiscordRed,
								Timestamp: discord.NowTimestamp(),
								Footer: &discord.EmbedFooter{
									Text: fmt.Sprintf("Ban issued by: %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
									Icon: e.Member.User.AvatarURL(),
								},
							},
						},
					},
				}

				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error.Println(err)
				}
			} else {
				res := api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Embeds: &[]discord.Embed{
							{
								Title:     "You can't ban yourself!!",
								Color:     utils.DiscordRed,
								Timestamp: discord.NowTimestamp(),
								Footer: &discord.EmbedFooter{
									Text: fmt.Sprintf("Ban attempt issued by: %v#%v", e.Member.User.Username, e.Member.User.Discriminator),
									Icon: e.Member.User.AvatarURL(),
								},
							},
						},
					},
				}

				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error.Println(err)
				}
			}
		},
	}
}
