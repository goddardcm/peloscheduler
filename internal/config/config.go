package config

import (
	"encoding/json"
	"os"
	"time"

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
	Peloton    Peloton    `json:"peloton"`
	OneMedical OneMedical `json:"one_medical"`
	OpenTable  OpenTable  `json:"open_table"`
	Twilio     Twilio     `json:"twilio"`
}

type Peloton struct {
	OrderID       string `json:"order_id"`
	QueryInterval string `json:"query_interval"`
}

type OneMedical struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	QueryInterval string `json:"query_interval"`
}

type OpenTable struct {
	RestaurantId  int64     `json:"restaurant_id"`
	NumberOfSeats int64     `json:"number_of_seats"`
	TimesToQuery  []string  `json:"times_to_query"`
	QueryDate     string    `json:"query_date"`
	LastDate      time.Time `json:"last_date"`
	QueryInterval string    `json:"query_interval"`
}

type Twilio struct {
	SID  string `json:"sid"`
	Auth string `json:"auth"`

	From string `json:"from"`
	To   string `json:"to"`
}
