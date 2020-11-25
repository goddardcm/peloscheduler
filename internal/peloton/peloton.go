package peloton

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const dateLayout = "2006-01-02T15:04:05-07:00"
const operationName = "OrderDelivery"
const requestQuery = `
query OrderDelivery($id: ID!, $isReschedule: Boolean = false) {
  order(id: $id) {
    canSetDeliveryPreference
    canReschedule
    canSelectTimeSlot
    deliveryPreference {
      date
      start
      end
      __typename
    }
    availableDeliveries(limit: 1, isReschedule: $isReschedule) {
      id
      date
      start
      end
      __typename
    }
    __typename
  }
  postPurchaseFlow(id: $id) {
    permission
    __typename
  }
}
`

type Order struct {
	CurrentDelivery     Date   `json:"deliveryPreference"`
	AvailableDeliveries []Date `json:"availableDeliveries"`
}

type Date struct {
	Date  string `json:"date"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func (d Date) GetStart() (time.Time, error) {
	return time.Parse(dateLayout, d.Start)
}

type graphqlRequest struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}
type graphqlResponse struct {
	Data struct {
		Order Order `json:"order"`
	} `json:"data"`
}

func FetchOrder(orderID string) (Order, error) {
	response := graphqlResponse{}

	reqBytes, reqErr := json.Marshal(
		graphqlRequest{
			OperationName: operationName,
			Variables: map[string]interface{}{
				"isReschedule": true,
				"id":           orderID,
			},
			Query: requestQuery,
		},
	)
	if reqErr != nil {
		return response.Data.Order, errors.WithStack(reqErr)
	}

	res, resErr := http.Post(
		"https://graph.prod.k8s.onepeloton.com/graphql",
		"application/json",
		bytes.NewReader(reqBytes),
	)
	if resErr != nil {
		return response.Data.Order, errors.WithStack(resErr)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(res.Body)
		return response.Data.Order, errors.Errorf("Received non-200 from Peloton: [%d] %s", res.StatusCode, string(responseBody))
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return response.Data.Order, errors.WithStack(err)
	}

	return response.Data.Order, nil
}
