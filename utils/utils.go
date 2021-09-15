package utils

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"log"
	"time"
)

// BToMb converts byte values to MegaBytes
func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
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

// the DefaultColour to use in embeds
var DefaultColour discord.Color = 0x7ECA9C

func MustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(env)
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}
