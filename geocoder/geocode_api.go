package geocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type GeocodeAPI struct {
	baseURL  string
	endpoint string
	apiKey   string
}

func NewGeoCodeMaps() Geocoder {
	return &GeocodeAPI{baseURL: "https://geocode.maps.co", endpoint: "/search", apiKey: os.Getenv("GEOCODE_API_KEY")}
}

func (g *GeocodeAPI) fullURL() string {
	return g.baseURL + g.endpoint
}

func (g *GeocodeAPI) getPlace(address string) ([]Place, error) {
	queryParams := url.Values{}
	queryParams.Add("q", address)
	queryParams.Add("api_key", g.apiKey)
	fullURL := fmt.Sprintf("%s?%s", g.fullURL(), queryParams.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		print("Error retrieving location information")
	}
	defer resp.Body.Close()

	var places []Place
	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	return places, nil
}
