package httputils

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func DoRequest(
	url string,
	authHeader string,
	requestStruct interface{},
	responsePtr interface{},
) error {
	reqBytes, reqErr := json.Marshal(requestStruct)
	if reqErr != nil {
		return errors.WithStack(reqErr)
	}

	httpRequest, httpRequestErr := http.NewRequest("POST", url, bytes.NewReader(reqBytes))
	if httpRequestErr != nil {
		return errors.WithStack(httpRequestErr)
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		httpRequest.Header.Set("Authorization", authHeader)
	}

	res, resErr := http.DefaultClient.Do(httpRequest)
	if resErr != nil {
		return errors.WithStack(resErr)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(res.Body)
		return errors.Errorf("Received non-200 from %w: [%d] %s", url, res.StatusCode, string(responseBody))
	}

	if err := json.NewDecoder(res.Body).Decode(responsePtr); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
