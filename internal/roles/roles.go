package roles

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kitokinha/terraform-provider/internal/client"
)

type Role struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Color       string  `json:"color"`
	IsDefault   bool    `json:"isDefault"`
	CreatedAt   string  `json:"createdAt"`
}

// List returns all default and custom roles visible to the configured token.
func List(c *client.PhaseClient) ([]Role, error) {
	endpoint := fmt.Sprintf("%s/v1/roles/", c.HostURL)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	responseBody, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var response struct {
		Data []Role `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}
