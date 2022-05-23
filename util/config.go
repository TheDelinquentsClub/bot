package util

import (
	"encoding/json"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// Config holds configuration values for the bot.
type Config struct {
	Path               string            `json:"-"`
	Token              string            `json:"token"`
	Version            string            `json:"version"`
	WolframID          string            `json:"wolfram_id"`
	OwnerID            discord.UserID    `json:"owner_id"`
	GuildID            discord.GuildID   `json:"guild_id"`
	MinecraftChannelID discord.ChannelID `json:"minecraft_channel_id"` // We fetch the appropriate discord.Webhook at runtime.
}

// NewConfig returns a new Config, reading from the provided file.
func NewConfig(filename string) (c *Config, err error) {
	if filename == "" {
		err = errors.New("no filename provided")
		return
	} else if !filepath.IsAbs(filename) {
		// We don't want to use filepath.Abs if we want the configuration file to be stored next to
		// the executable.
		self, err2 := os.Executable()
		if err2 != nil {
			err = err2
			return
		}
		self, err2 = filepath.EvalSymlinks(self)
		if err2 != nil {
			err = err2
			return
		}
		filename = filepath.Join(filepath.Dir(self), filename)
	}

	c = &Config{Path: filename}
	if err = c.Update(); err != nil {
		c = nil
		return
	}

	return
}

// Update updates the Config from its already-defined Path. If you want to change its path, reassign
// and call this afterwards.
func (c *Config) Update() (err error) {
	if c.Path == "" {
		return errors.New("(*Config).Path cannot be empty")
	}
	file, err2 := os.Open(c.Path)
	if err2 != nil {
		err = err2
		return
	}
	// We don't use defer here, because we can avoid it.

	temp := Config{Path: c.Path}
	if err = json.NewDecoder(file).Decode(&temp); err != nil {
		return
	} else if err = temp.validate(); err != nil {
		return
	}
	*c, _ = temp, file.Close()
	return nil
}

func (c *Config) validate() error {
	switch {
	case c.Token == "" || len(c.Token) < 5:
		return errors.New("(*Config).Token cannot be empty or smaller than 5 characters")
	case c.Token[:5] != "Bot ":
		if c.Token[:4] == "Bot" {
			c.Token = "Bot " + c.Token[4:]
		} else {
			c.Token = "Bot " + c.Token
		}
	case c.Version == "":
		c.Version = "Unknown"
	case c.WolframID == "":
		return errors.New("(*Config).WolframID cannot be empty")
	case !c.OwnerID.IsValid():
		return errors.New("(*Config).OwnerID must be valid")
	case !c.GuildID.IsValid():
		return errors.New("(*Config).GuildID must be valid")
	case !c.MinecraftChannelID.IsValid():
		return errors.New("(*Config).MinecraftChannelID must be valid")
	}
	return nil
}
