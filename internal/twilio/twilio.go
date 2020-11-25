package twilio

import (
	"github.com/goddardcm/peloscheduler/internal/config"
	"github.com/kevinburke/twilio-go"
	"github.com/pkg/errors"
)

func SendMessage(config config.Twilio, message string) error {
	_, sendErr := twilio.NewClient(config.SID, config.Auth, nil).
		Messages.
		SendMessage(config.From, config.To, message, nil)

	return errors.WithStack(sendErr)
}
