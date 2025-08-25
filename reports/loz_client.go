package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseUrl = "https://botx-compendium.herokuapp.com/api/v3/compendium"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LozClient struct {
	baseUrl    string
	httpClient HttpClient
}

func NewLozClient(httpClient HttpClient) *LozClient {
	return &LozClient{
		baseUrl:    baseUrl,
		httpClient: httpClient,
	}
}

type Monster struct {
	Name            string   `json:"name"`
	Id              int      `json:"id"`
	Category        string   `json:"category"`
	Description     string   `json:"description"`
	Image           string   `json:"image"`
	CommonLocations []string `json:"common_locations"`
	Drops           []string `json:"drops"`
	Dlc             bool     `json:"dlc"`
}

type GetMonsterResponse struct {
	Data []Monster `json:"data"`
}

func (c *LozClient) GetMonsters() (*GetMonsterResponse, error) {
	req, err := http.NewRequest("GET", c.baseUrl+"/category/monsters", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create monsters request: %w", err)
	}

	reqUrl := req.URL
	queryParams := req.URL.Query()
	queryParams.Set("game", "totk")
	reqUrl.RawQuery = queryParams.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get monsters http request: %w", err)
	}

	var response *GetMonsterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to get monsters http response: %w", err)
	}

	return response, nil

}
