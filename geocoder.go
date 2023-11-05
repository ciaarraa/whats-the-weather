package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://geocode.maps.co"
const endpoint = "/search"

type HTTPClient struct {
	baseURL string
}

func main() {

	type Place struct {
		Latitude    string `json:"lat"`
		Longitude   string `json:"lon"`
		DisplayName string `json:"display_name"`
	}

	queryParams := url.Values{}
	queryParams.Add("q", "Statue of Liberty NY US")
	fullURL := fmt.Sprintf("%s?%s", baseURL+endpoint, queryParams.Encode())
	fmt.Println(fullURL)
	resp, _ := http.Get(fullURL)
	defer resp.Body.Close()

	var places []Place

	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, place := range places {
		fmt.Println("Display Name:", place.DisplayName)
		fmt.Println("Latitude:", place.Latitude)
		fmt.Println("Longitude", place.Longitude)
	}
}
