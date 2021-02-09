package onemedical

import (
	"fmt"
	"github.com/goddardcm/peloscheduler/internal/config"
	"github.com/goddardcm/peloscheduler/internal/httputils"
	"time"
)

const dateLayout = "2006-01-02"

func FetchAppointments(config config.OneMedical) (string, error) {
	accessToken, accessTokenErr := fetchAccessToken(config)
	if accessTokenErr != nil {
		return "", accessTokenErr
	}

	return fetchAppointmentAvailable(accessToken)
}

func fetchAppointmentAvailable(accessToken string) (string, error) {
	request := appointmentRequest{
		AppointmentTypeID: 275,
		ServiceAreaID:     1,
		OfficeIDs:         make([]int64, 0),
		ProviderID:        nil,
		OnsiteOnly:        false,

		StartDate: time.Now().Format(dateLayout),
		EndDate:   time.Now().Add(2 * 24 * time.Hour).Format(dateLayout),
	}
	response := appointmentResponse{}

	httpErr := httputils.DoRequest(
		"https://members.onemedical.com/api/v2/patient/appointment_search",
		fmt.Sprintf("Bearer %s", accessToken),
		request,
		&response,
	)

	return response.FirstAvailableInventoryDate, httpErr
}

func fetchAccessToken(config config.OneMedical) (string, error) {
	request := authRequest{
		Username:  config.Username,
		Password:  config.Password,
		GrantType: "password",
	}
	response := authResponse{}

	httpErr := httputils.DoRequest(
		"https://members.onemedical.com/oauth/token",
		"",
		request,
		&response,
	)

	return response.AccessToken, httpErr
}
