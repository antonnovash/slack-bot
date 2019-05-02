package bot

import (
	"context"
	"fmt"
	"log"
	"slack-bot/slack-bot/pkg/slack"

	slackapi "github.com/nlopes/slack"
)

// Bot contains managed websocket connection.
type Bot struct {
	rtm   *slackapi.RTM
	slack *slack.SlackListener
}

// New creates an instance of Slack bot.
// If succeed returns new instance and nil.
func New(c slack.Config, slack *slack.SlackListener) (*Bot, error) {
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config is invalid: %v", err)
	}
	client := slackapi.New(c.Token)
	rtm := client.NewRTM()
	return &Bot{rtm: rtm, slack: slack}, nil
}

// Run starts handling of incoming RTMEvents.
func (b *Bot) Run(ctx context.Context) error {
	go b.rtm.ManageConnection()
	for {
		select {
		case incomingEvent, ok := <-b.rtm.IncomingEvents:
			if !ok {
				return fmt.Errorf("channel IncomingEvents was closed")
			}
			switch event := incomingEvent.Data.(type) {
			case *slackapi.MessageEvent:
				_ = b.slack.HandleMessageEvent(event)
			case *slackapi.RTMError:
				log.Printf("slack runtime error: %s, %d, %v", event.Msg, event.Code, event.Error())
			case *slackapi.InvalidAuthEvent:
				return fmt.Errorf("invalid credentials")
			}
		case <-ctx.Done():
			err := b.rtm.Disconnect()
			if err != nil {
				return fmt.Errorf("cannot disconnect successfully: %v", err)
			}
			return nil
		}
	}
}
