// reservations-api/clients_hotels.go
package clients_reservations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HotelsClient interface {
	HotelExists(id string) bool
}

type HTTPHotelsClient struct {
	BaseURL string
	Client  *http.Client
}

func NewHTTPHotelsClient(baseURL string) *HTTPHotelsClient {
	return &HTTPHotelsClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

type hotelDTO struct {
	ID string `json:"id"`
}

func (c *HTTPHotelsClient) HotelExists(id string) bool {
	if id == "" {
		return false
	}
	url := fmt.Sprintf("%s/hotels/%s", c.BaseURL, id)
	resp, err := c.Client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	var h hotelDTO
	if err := json.NewDecoder(resp.Body).Decode(&h); err != nil {
		return false
	}
	return h.ID != ""
}
