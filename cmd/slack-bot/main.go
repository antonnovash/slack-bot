package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slack-bot/slack-bot/pkg/bot/slack"
	"slack-bot/slack-bot/pkg/config"
	"slack-bot/slack-bot/pkg/server"
	"slack-bot/slack-bot/pkg/slackApp"
	"syscall"
)

// Config stores info from env vars.
type Config struct {
	Slack  slackApp.Config
	Server server.Config
}

func main() {
	c := new(config.Config)
	if err := c.Load(config.SERVICENAME); err != nil {
		log.Fatal(err)
	}
	//TODO Add Exchange,Outlook,DB connection
	s, err := customize(c)
	err = ServerRun(c)
	fmt.Println(c)
	if err != nil {
		log.Fatal(err)
	}
	s.BotRun()
}
func ServerRun(c *config.Config) error {
	myChan := make(chan string)
	server, err := server.New(c.ServerAddress, myChan)
	if err != nil {
		return fmt.Errorf("cannot connect to Server %v", err)
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
	log.Println("successful token received")
	c.BotToken = token
	return nil
}


func customize(c *config.Config) (*slack.SlackBot, error) {
	slackApps, err := slackApp.New(c)
	if err != nil {
		return nil, fmt.Errorf("cannot create slack connection: %v", err)
	}
	bot, err := slack.New(c, slackApps)
	if err != nil {
		return nil, fmt.Errorf("cannot create a bot: %v", err)
	}
	return bot, err
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