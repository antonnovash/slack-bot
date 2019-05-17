package slack

import "fmt"

// Config contains info Slack bot needs.
type Config struct {
	Token     string
	ChannelID string
}

// Validate validates Slack bot config.
func (c Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("slack bot token env variable isn't set: %s", c.Token)
	}
	if c.ChannelID == "" {
		return fmt.Errorf("slack channel ID env variable isn't set: %s", c.ChannelID)
	}
	return nil
}
