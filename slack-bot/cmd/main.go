package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slack-bot/slack-bot/pkg/bot"
	"slack-bot/slack-bot/pkg/controller"
	"slack-bot/slack-bot/pkg/server"
	"slack-bot/slack-bot/pkg/slack"
	"syscall"
)

// Environment variables.
const (
	envSlackChannelID = "SLACK_CHANNEL"
	envSlackToken     = "SLACK_TOKEN"
	envServerAddress  = "SERVER_ADDRESS"
)
// Config stores info from env vars.
type Config struct {
	Slack  slack.Config
	Server server.Config
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
	slack, err := slack.New(config.Slack)
	if err != nil {
		return fmt.Errorf("cannot create slack connection: %v", err)
	}

	bot, err := bot.New(config.Slack, slack)
	if err != nil {
		return fmt.Errorf("cannot create a bot: %v", err)
	}
	controller := &controller.Controller{
		Slack: slack,
	}
	myChan := make(chan string)

	server, err := server.New(config.Server, controller, myChan)
	if err != nil {
		return fmt.Errorf("cannot create a server: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	shutdownOnSignal(cancel)
	errChan := make(chan error, 1)
	go func() {
		//config.Slack.Token =  <-myChan
		token :=  <-myChan
		fmt.Println(token)
		err = bot.Run(ctx)
				if err != nil {
			cancel()
			errChan <- fmt.Errorf("bot stopped with error: %v", err)
			return
		}
		errChan <- nil
	}()

	go func() {
		err := server.Run(ctx)
		if err != nil {
			cancel()
			errChan <- fmt.Errorf("server stopped with error: %v", err)
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

// NewConfig returns a Config struct.
func NewConfig() (*Config, error) {
	os.Setenv(envSlackChannelID, "CJ3UZQ6P7")
	os.Setenv(envServerAddress, "localhost:9000")
	c := &Config{

		Slack: slack.Config{
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
