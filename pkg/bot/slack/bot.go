package slack

import (
	"context"
	"fmt"
	slackapi "github.com/nlopes/slack"
	"log"
	"os"
	"os/signal"
	"slack-bot/slack-bot/pkg/config"
	"slack-bot/slack-bot/pkg/slackApp"
	"syscall"
	"time"
)

type command int

const (
	noneCommand command = iota
	tokenCommand
	subjectCommand
	emailCommand
)

type meetingRequest struct {
	StartTime time.Time
	Duration  time.Duration
	/*AvailableRooms []entity.Room
	ChosenRoom     entity.Room*/
}

// Bot contains managed websocket connection.
type SlackBot struct {
	rtm   *slackapi.RTM
	slack *slackApp.SlackListener
	//	ctrl  ctrl.Ctrl TODO add Controller
	users map[int]*slackUser
}

type slackUser struct {
	meetingRequest meetingRequest
	command        command
}

// New creates an instance of Slack bot.
// If succeed returns new instance and nil.
func New(c *config.Config, slack *slackApp.SlackListener) (*SlackBot, error) {
	//TODO add validation
	client := slackapi.New(c.BotToken)
	rtm := client.NewRTM()
	return &SlackBot{rtm: rtm, slack: slack}, nil
}

func (s *SlackBot) BotRun() {
	ctx, cancel := context.WithCancel(context.Background())
	err := s.Run(ctx)
	if err != nil {
		cancel()
		fmt.Errorf("bot stopped with error: %v", err)
		return
	}
}

// Run starts handling of incoming RTMEvents.
func (s *SlackBot) Run(ctx context.Context) error {
	go s.rtm.ManageConnection()
	for {
		select {
		case incomingEvent, ok := <-s.rtm.IncomingEvents:
			if !ok {
				return fmt.Errorf("channel IncomingEvents was closed")
			}
			switch event := incomingEvent.Data.(type) {
			case *slackapi.MessageEvent:
				_ = s.slack.HandleMessageEvent(event)
			case *slackapi.RTMError:
				log.Printf("slack runtime error: %s, %d, %v", event.Msg, event.Code, event.Error())
			case *slackapi.InvalidAuthEvent:
				return fmt.Errorf("invalid credentials")
			}
		case <-ctx.Done():
			err := s.rtm.Disconnect()
			if err != nil {
				return fmt.Errorf("cannot disconnect successfully: %v", err)
			}
			return nil
		}
	}
}
func shutdownOnSignal(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-stop
		cancel()
		log.Printf("Shutting down due to signal: %v", sig)
	}()
}
