package main

import (
	"SlackBot/slack-bot/pkg/bot"
	"SlackBot/slack-bot/pkg/slack"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Environment variables.
const (
	envSlackChannelID = "SLACK_CHANNEL"
	envSlackToken     = "SLACK_TOKEN"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = run(config)
	if err != nil {
		log.Fatal(err)
	}
}

func run(config *Config) error {
	slack, err := slack.New(config.Slack)
	if err != nil {
		return fmt.Errorf("cannot create slack connection: %v", err)
	}

	bot, err := bot.New(config.Slack, slack)
	if err != nil {
		return fmt.Errorf("cannot create a bot: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	shutdownOnSignal(cancel)
	errChan := make(chan error, 1)
/*	go func() {
		if err != nil {
			cancel()
			errChan <- fmt.Errorf("server stopped with error: %v", err)
			return
		}
		errChan <- nil

	}()*/
	go func() {
		err = bot.Run(ctx)
		if err != nil {
			cancel()
			errChan <- fmt.Errorf("bot stopped with error: %v", err)
			return
		}
		errChan <- nil
	}()
	err = <-errChan
	<-errChan
	return err
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


// Config stores info from env vars.
type Config struct {
	Slack slack.Config
}
// NewConfig returns a Config struct.
func NewConfig() (*Config, error) {
	os.Setenv(envSlackToken, "xoxb-617863072727-610208349666-j8yCGWNWpyTWfPK8r5p1atu6")
	os.Setenv(envSlackChannelID, "CJ3UZQ6P7")
	c := &Config{

		Slack: slack.Config{
			Token:     os.Getenv(envSlackToken),
			ChannelID: os.Getenv(envSlackChannelID),
		},
	}
	err := c.validate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c Config) validate() error {
	err := c.Slack.Validate()
	if err != nil {
		return fmt.Errorf("bot config: %v", err)
	}
	return nil
}
