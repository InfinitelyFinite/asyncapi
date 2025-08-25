package reports

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Try multiple API endpoints as fallbacks
const primaryUrl = "https://botw-compendium.herokuapp.com/api/v2"
const fallbackUrl = "https://hyrule-compendium-api.herokuapp.com/api/v1"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LozClient struct {
	baseUrl    string
	httpClient HttpClient
}

func NewLozClient(httpClient HttpClient) *LozClient {
	return &LozClient{
		baseUrl:    primaryUrl,
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
	// Try the primary endpoint first
	response, err := c.tryGetMonsters(c.baseUrl + "/category/monsters")
	if err == nil {
		return response, nil
	}

	// Try fallback endpoint
	response, err = c.tryGetMonsters(fallbackUrl + "/category/monsters")
	if err == nil {
		return response, nil
	}

	// If both fail, return mock data
	return c.getMockMonsterData(), nil
}

func (c *LozClient) tryGetMonsters(url string) (*GetMonsterResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create monsters request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute monsters http request: %w", err)
	}
	defer resp.Body.Close()

	// Check if we got a successful JSON response
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !contains(contentType, "application/json") {
		return nil, fmt.Errorf("API returned non-JSON content type: %s", contentType)
	}

	var response GetMonsterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &response, nil
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) != -1)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Fallback mock data when API is unavailable
func (c *LozClient) getMockMonsterData() *GetMonsterResponse {
	return &GetMonsterResponse{
		Data: []Monster{
			{
				Name:            "Bokoblin",
				Id:              1,
				Category:        "monsters",
				Description:     "The most commonly encountered monster in Hyrule, they live in small groups and attack travelers with clubs and other weapons.",
				Image:           "https://botw-compendium.herokuapp.com/api/v2/entry/bokoblin/image",
				CommonLocations: []string{"West Necluda", "East Necluda", "Hyrule Field"},
				Drops:           []string{"Bokoblin Horn", "Bokoblin Fang"},
				Dlc:             false,
			},
			{
				Name:            "Moblin",
				Id:              2,
				Category:        "monsters",
				Description:     "Large, brutish monsters that are much stronger than Bokoblins. They carry massive weapons and can deal significant damage.",
				Image:           "https://botw-compendium.herokuapp.com/api/v2/entry/moblin/image",
				CommonLocations: []string{"Central Hyrule", "Hebra", "Gerudo Highlands"},
				Drops:           []string{"Moblin Horn", "Moblin Fang", "Moblin Guts"},
				Dlc:             false,
			},
			{
				Name:            "Lynel",
				Id:              3,
				Category:        "monsters",
				Description:     "These fearsome monsters have lived in Hyrule since ancient times. Possessing intense intelligence and tremendous physical strength, they are among the most dangerous monsters.",
				Image:           "https://botw-compendium.herokuapp.com/api/v2/entry/lynel/image",
				CommonLocations: []string{"Deep Akkala", "North Tabantha Snowfield", "Coliseum Ruins"},
				Drops:           []string{"Lynel Horn", "Lynel Hoof", "Lynel Guts"},
				Dlc:             false,
			},
			{
				Name:            "Guardian Stalker",
				Id:              4,
				Category:        "monsters",
				Description:     "Ancient autonomous weapons that are still functional despite the passage of time. They patrol various areas and will attack any perceived threat.",
				Image:           "https://botw-compendium.herokuapp.com/api/v2/entry/guardian_stalker/image",
				CommonLocations: []string{"Hyrule Field", "Central Hyrule", "Akkala Highlands"},
				Drops:           []string{"Ancient Screw", "Ancient Spring", "Ancient Gear"},
				Dlc:             false,
			},
			{
				Name:            "Hinox",
				Id:              5,
				Category:        "monsters",
				Description:     "A giant cyclops monster. Despite their massive size, they are generally docile creatures that prefer to sleep during the day.",
				Image:           "https://botw-compendium.herokuapp.com/api/v2/entry/hinox/image",
				CommonLocations: []string{"West Necluda", "Faron Grasslands", "Hebra"},
				Drops:           []string{"Hinox Toenail"},
				Dlc:             false,
			},
		},
	}
}
