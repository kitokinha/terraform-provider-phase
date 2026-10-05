package serviceaccounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kitokinha/terraform-provider/internal/client"
)

func Create(c *client.PhaseClient, account CreateRequest) (*ServiceAccount, error) {
	return request[CreateRequest, ServiceAccount](c, http.MethodPost, "/v1/service-accounts/", account)
}

func Read(c *client.PhaseClient, id string) (*ServiceAccount, error) {
	url := fmt.Sprintf("/v1/service-accounts/%s/", id)
	return request[any, ServiceAccount](c, http.MethodGet, url, nil)
}

func List(c *client.PhaseClient) ([]ServiceAccount, error) {
	url := fmt.Sprintf("%s/v1/service-accounts/", c.HostURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	responseBody, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var response struct {
		Data []ServiceAccount `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func Update(c *client.PhaseClient, id string, account UpdateRequest) (*ServiceAccount, error) {
	url := fmt.Sprintf("/v1/service-accounts/%s/", id)
	return request[UpdateRequest, ServiceAccount](c, http.MethodPut, url, account)
}

// SetAccess replaces the service account's complete app/environment access set.
func SetAccess(c *client.PhaseClient, id string, access AccessRequest) (*ServiceAccount, error) {
	url := fmt.Sprintf("/v1/service-accounts/%s/access/", id)
	return request[AccessRequest, ServiceAccount](c, http.MethodPut, url, access)
}

func Delete(c *client.PhaseClient, id string) error {
	url := fmt.Sprintf("%s/v1/service-accounts/%s/", c.HostURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func request[Request any, Response any](c *client.PhaseClient, method, path string, body Request) (*Response, error) {
	var requestBody *bytes.Buffer
	if method == http.MethodGet {
		requestBody = bytes.NewBuffer(nil)
	} else {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(encoded)
	}

	req, err := http.NewRequest(method, c.HostURL+path, requestBody)
	if err != nil {
		return nil, err
	}
	responseBody, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
