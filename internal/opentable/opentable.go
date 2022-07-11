package opentable

import (
	"fmt"
	"time"

	"github.com/goddardcm/peloscheduler/internal/config"
	"github.com/goddardcm/peloscheduler/internal/httputils"
	"github.com/pkg/errors"
)

const urlFormat = "https://www.opentable.com/restaurant/profile/%d/search"
const dateLayout = "2006-01-02"

type searchRequest struct {
	Covers        int64  `json:"covers"`
	Time          string `json:"dateTime"`
	IsRedesign    bool   `json:"isRedesign"`
	CorrelationId string `json:"correlationId"`
}

type availability struct {
	Date  string `json:"date"`
	Times []struct {
		Time     string    `json:"timeString"`
		DateTime time.Time `json:"dateTime"`
	} `json:"times"`
}

type searchResponse struct {
	MultiDaysAvailability struct {
		Times []availability `json:"timeslots"`
	} `json:"multiDaysAvailability"`
}

func SearchForAvailability(config config.OpenTable) ([]string, error) {
	returnValue := make([]string, 0)

	url := fmt.Sprintf(urlFormat, config.RestaurantId)
	for _, queryTime := range config.TimesToQuery {
		request := searchRequest{
			Covers:     config.NumberOfSeats,
			Time:       fmt.Sprintf("%sT%s", config.QueryDate, queryTime),
			IsRedesign: true,
		}
		response := searchResponse{}

		if httpErr := httputils.DoRequest(
			url,
			"",
			request,
			&response,
		); httpErr != nil {
			return returnValue, errors.Wrap(httpErr, "Error querying OpenTable")
		}

		for _, availability := range response.MultiDaysAvailability.Times {
			for _, availabilityTime := range availability.Times {
				if availabilityTime.DateTime.After(config.LastDate) {
					continue
				}
				returnValue = append(returnValue, fmt.Sprintf("%s @ %s", availability.Date, availabilityTime.Time))
			}
		}
	}

	return returnValue, nil

}
