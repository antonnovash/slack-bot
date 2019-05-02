package server

import "fmt"

// Config contains server's address.
type Config struct {
	Address string
}

// Validate validates server config.
func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("env SERVER_ADDRESS isn't set")
	}
	return nil
}
