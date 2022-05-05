package utils

import (
	"encoding/json"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/logger"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func ConvertTickToDuration(ticks int) string {
	var years = ticks / 630720000
	ticks = ticks % 630720000
	var days = ticks / 1728000
	ticks = ticks % 1728000
	var hours = ticks / 72000
	ticks = ticks % 72000
	var minutes = ticks / 1200
	ticks = ticks % 1200
	var seconds = ticks / 20
	ticks = ticks % 20
	return fmt.Sprintf("%v:%v:%v:%v:%v", years, days, hours, minutes, seconds)
}

// DefaultColour is the default discord.color to use in embeds
// DiscordGreen is the colour to be used in signifying a success message, or something good
// DiscordRed is the colour to be used in signifying an error message, or something bad
var (
	DefaultColour discord.Color = 0xA3BCF9
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

type Player struct {
	Username string `json:"name"`
	UUID     string `json:"id"`
}

// GetUUID returns the UUID tied with to the username provided
func GetUUID(username string) string {
	username = strings.ReplaceAll(username, "\"", "")
	url := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%v", username)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	var user = new(Player)
	err = json.Unmarshal(body, &user)
	return user.UUID

}

// GetUsername returns the username of a player from the provided UUID
func GetUsername(UUID string) string {
	UUID = strings.ReplaceAll(UUID, "\"", "")
	url := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%v", UUID)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	var user = new(Player)
	err = json.Unmarshal(body, &user)
	return user.Username
}

type PlayerNames struct {
	Name    string `json:"name"`
	Changed int64  `json:"changedToAt,omitempty"`
}

// GetNamesFromUsername returns all the usernames the specified player has had using the provided username
func GetNamesFromUsername(username string) string {
	uuid := GetUUID(username)
	return GetNamesFromUUID(uuid)
}

// GetNamesFromUUID returns all the usernames the specified player has had using the provided UUID
func GetNamesFromUUID(uuid string) string {
	uuid = strings.ReplaceAll(uuid, "\"", "")
	url := fmt.Sprintf("https://api.mojang.com/user/profiles/%v/names", uuid)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	var names = new([]PlayerNames)
	err = json.Unmarshal(body, &names)

	res := strings.Builder{}

	for _, playerNames := range *names {
		unixString := strings.TrimSuffix(fmt.Sprint(playerNames.Changed), "000")
		unixTime, err := strconv.ParseInt(unixString, 10, 64)
		if err != nil {
			logger.Error(err)
		}

		changed := fmt.Sprintf("Changed: <t:%v:R>\n\n", time.Unix(unixTime, 0))
		if playerNames.Changed == 0 {
			changed = "Accounts first username!\n\n"
		}
		res.WriteString(fmt.Sprintf("Username: %v\n%v", playerNames.Name, changed))
	}

	return res.String()
}
