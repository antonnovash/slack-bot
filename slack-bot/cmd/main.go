package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slack-bot/slack-bot/pkg/bot"
	"slack-bot/slack-bot/pkg/server"
	"slack-bot/slack-bot/pkg/slack"
	"syscall"
)

// Environment variables.
const (
	envSlackChannelID    = "SLACK_CHANNEL"
	envSlackToken        = "SLACK_TOKEN"
	envServerAddress     = "SERVER_ADDRESS"
	envSlackClientId     = "SLACK_CLIENT_ID"
	envSlackClientSecret = "SLACK_CLIENT_SECRET"
)

// Config stores info from env vars.
type Config struct {
	Slack  slack.Config
	Server server.Config
}

// NewConfig returns a Config struct.
func NewConfig() (*Config, error) {
	_ = os.Setenv(envSlackClientId, "617863072727.609868118610")
	_ = os.Setenv(envSlackClientSecret, "a39e1d00bbe6ce9a88c191391108600c")
	_ = os.Setenv(envSlackToken, "xoxp-617863072727-604564212419-623675919763-604c125b972fba3e9d5e40067dbe4557")
	_ = os.Setenv(envSlackChannelID, "CJ3UZQ6P7")
	_ = os.Setenv(envServerAddress, "localhost:9000")
	c := &Config{

		Slack: slack.Config{
			Token:     os.Getenv(envSlackToken),
			ChannelID: os.Getenv(envSlackChannelID),
		},
		Server: server.Config{
			Address: os.Getenv(envServerAddress),
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
	myChan := make(chan string)
	server, err := server.New(config.Server, myChan)
	if err != nil {
		return fmt.Errorf("cannot create a server: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	shutdownOnSignal(cancel)
	errChan := make(chan error, 1)
	go func() {
		err := server.Run(ctx)
		if err != nil {
			cancel()
			errChan <- fmt.Errorf("server stopped with error: %v", err)
			return
		}
		errChan <- nil
	}()
	token := <-myChan
	fmt.Println("token get sucsess.")
	config.Slack.Token = token
	slack, err := slack.New(config.Slack)
	if err != nil {
		return fmt.Errorf("cannot create slack connection: %v", err)
	}
	bot, err := bot.New(config.Slack, slack)
	if err != nil {
		return fmt.Errorf("cannot create a bot: %v", err)
	}
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
	fmt.Println(err)
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
