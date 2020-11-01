package script

import "github.com/pkg/errors"

func DefaultConfig() *Config {
	return &Config{
		Location: ".",
	}
}

type Config struct {
	Location string
}

func (c *Config) Validate() error {
	if len(c.Location) == 0 {
		return errors.New("location cannot be empty")
	}

	return nil
}
