package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

func FromFile(filePath string) (Config, error) {
	config := Config{}

	file, openErr := os.Open(filePath)
	if openErr != nil {
		return config, errors.WithStack(openErr)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, errors.WithStack(err)
	}

	return config, nil
}

type Config struct {
	Peloton Peloton `json:"peloton"`
	Twilio  Twilio  `json:"twilio"`
}

type Peloton struct {
	OrderID       string `json:"order_id"`
	QueryInterval string `json:"query_interval"`
}

type Twilio struct {
	SID  string `json:"sid"`
	Auth string `json:"auth"`

	From string `json:"from"`
	To   string `json:"to"`
}
