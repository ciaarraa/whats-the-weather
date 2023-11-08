package geocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://geocode.maps.co"
const endpoint = "/search"

type GeocoderClient struct{}

type Place struct {
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	DisplayName string `json:"display_name"`
}

func NewGeocoderClient() *GeocoderClient {
	return &GeocoderClient{}
}

func (g *GeocoderClient) FindCoordinates(address string) ([]Place, error) {
	queryParams := url.Values{}
	queryParams.Add("q", address)
	fullURL := fmt.Sprintf("%s?%s", baseURL+endpoint, queryParams.Encode())
	resp, _ := http.Get(fullURL)

	var places []Place

	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	defer resp.Body.Close()
	// TO DO: Add proper error response here
	return places, nil
}
