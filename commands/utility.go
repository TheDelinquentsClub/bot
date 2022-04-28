package commands

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/kingultron99/tdcbot/core"
	"github.com/kingultron99/tdcbot/logger"
	"github.com/kingultron99/tdcbot/utils"
)

func init() {
	MapCommands["ban"] = Command{
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
					logger.Error(fmt.Sprintf("Failed to ban user %v\nerror: %v", user.Username, err))
				}

				res := utils.NewEmbed().
					SetTitle("🔨 Banned!").
					AddField("User", false, fmt.Sprintf("<@%v>", user.ID)).
					AddField("Reason", true, reason).
					SetColor(utils.DiscordRed).
					SetFooter(fmt.Sprintf("Ban issued by: %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
					MakeResponse()

				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
			} else {
				res := utils.NewEmbed().
					SetTitle("You can't ban yourself!!").
					SetColor(utils.DiscordRed).
					SetFooter(fmt.Sprintf("Ban attempt issued by: %v#%v", e.Member.User.Username, e.Member.User.Discriminator), e.Member.User.AvatarURL()).
					MakeResponse()

				if err := core.State.RespondInteraction(e.ID, e.Token, res); err != nil {
					logger.Error(err)
				}
			}
		},
	}
}
