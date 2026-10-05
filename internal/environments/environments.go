package environments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kitokinha/terraform-provider/internal/client"
)

type Environment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	EnvType   string `json:"envType"`
	Index     int    `json:"index"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// List returns the environments accessible to the configured token for an app.
func List(c *client.PhaseClient, appID string) ([]Environment, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/v1/environments/", c.HostURL))
	if err != nil {
		return nil, err
	}
	query := endpoint.Query()
	query.Set("app_id", appID)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	responseBody, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var response struct {
		Data []Environment `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}
