package controller

import (
	"fmt"
	"slack-bot/slack-bot/pkg/slack"
)

// Controller contains connections to Exchange, YouTube and Slack bot.
type Controller struct {
	Slack    *slack.SlackListener
}

// Validate checks if all field are not nil.
func (c Controller) Validate() error {
	if c.Slack == nil {
		return fmt.Errorf("slack is nil")
	}
	return nil
}
