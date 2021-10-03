package utils

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"log"
	"strings"
	"time"
)

// BToMb converts byte values to MegaBytes
func BToMb(b uint64) string {
	return fmt.Sprintf("%.1f", float64(b/1024)/1024)
}

// GetDurationString returns a time duration for use with uptime or other related uses
func GetDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

// DefaultColour is the default discord.color to use in embeds
// DiscordGreen is the colour to be used in signifying a success message, or something good
// DiscordRed is the colour to be used in signifying an error message, or something bad
var (
	DefaultColour discord.Color = 0x7ECA9C
	DiscordGreen  discord.Color = 0x379A57
	DiscordBlue   discord.Color = 0x5865F2
	DiscordRed    discord.Color = 0xDF3E41
)

func MustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(env)
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}

func GenButtonComponents(data discord.InteractionOption) *discord.ButtonComponent {
	component := &discord.ButtonComponent{
		Label:    strings.ReplaceAll(fmt.Sprint(data.Value), "\"", ""),
		CustomID: fmt.Sprint(data.Name),
		Style:    discord.PrimaryButton,
	}

	return component
}
