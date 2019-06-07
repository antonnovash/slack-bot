package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	// SERVICENAME contains a service name prefix which used in ENV variables
	SERVICENAME = "SLACKBOT"
)

// Config contains ENV variables
type Config struct {
	// Exchange username
	ExchangeUser string `split_words:"true" required:"true"`
	// Exchange password
	ExchangePass string `split_words:"true" required:"true"`
	// Connection string for database
	ConnectionString string `split_words:"true" required:"true"`
	// ExchangeURL is an URL that used to connect to exchange server
	ExchangeURL string `split_words:"true" required:"true"`
	// Slack bot api token
	BotToken string `split_words:"true" required:"true"`
	//SlackChannelID
	ChannelID string `split_words:"true" required:"true"`
	//ServerAddress
	ServerAddress string `split_words:"true" required:"true"`
	//SlackClientId
	ClientID string `split_words:"true" required:"true"`
	//SlackClientSecret
	ClientSecret string `split_words:"true" required:"true"`

	ControllerConfig
}

// ControllerConfig contains variables, that are needed for some of the controller functions
type ControllerConfig struct {
	// BotName is a bot name that books rooms
	BotName string `split_words:"true" required:"true"`
	// RoomBlackList is list of rooms, that cannot be chosen. Rooms are written as emails, and divided with `,`
	RoomBlacklist string `split_words:"true"`
	// ExchangeNewMeetingDaysLimit is a limit of the number of days from now in which a meeting can be created.
	ExchangeNewMeetingDaysLimit int `split_words:"true" default:"120"`
}

// Load settles ENV variables into Config structure
func (c *Config) Load(serviceName string) error {
	return envconfig.Process(serviceName, c)
}
