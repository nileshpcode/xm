package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	apiError "xm/error"
)

type IPLocationClient interface {
	GetLocation(ip string) (string, error)
}

// NewIpLocationClient returns a new instance of IpLocationClient
func NewIpLocationClient(url string) IPLocationClient {
	client := &ipLocationClientImpl{}
	client.HTTPClient = &http.Client{}
	client.BaseURL = url
	return client
}

type ipLocationClientImpl struct {
	BaseURL    string
	HTTPClient *http.Client
}

// GetLocation gets location of ip
func (impl *ipLocationClientImpl) GetLocation(ip string) (string, error) {
	apiURL := fmt.Sprintf("%s/%s/json/", impl.BaseURL, ip)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", apiError.NewAPIClientError(apiURL, nil, nil, fmt.Errorf("unable to create HTTP request: %w", err))
	}

	req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")

	resp, err := impl.HTTPClient.Do(req)
	if err != nil {
		return "", apiError.NewAPIClientError(apiURL, nil, nil, fmt.Errorf("unable to invoke API: %w", err))
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		responseBodyString := ""
		if responseBodyBytes, err := ioutil.ReadAll(resp.Body); err == nil {
			responseBodyString = string(responseBodyBytes)
		}
		return "", apiError.NewAPIClientError(apiURL, &resp.StatusCode, &responseBodyString, fmt.Errorf("received non-ok code: %v", resp.StatusCode))
	}

	var mapResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mapResponse); err != nil {
		return "", apiError.NewAPIClientError(apiURL, &resp.StatusCode, nil, fmt.Errorf("unable parse response payload: %w", err))
	}

	if isErr, ok := mapResponse["error"].(bool); ok && isErr {
		return "", apiError.NewAPIClientError(apiURL, nil, nil, fmt.Errorf(mapResponse["reason"].(string)))
	}

	return mapResponse["country"].(string), nil
}
