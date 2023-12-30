package geocoder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Geocoder interface {
	getPlace(geocodeURL string) ([]Place, error)
	fullURL() string
}

type GeocodeAPI struct {
	baseURL  string
	endpoint string
	apiKey   string
}

const defualtBaseURL = "https://geocode.maps.co"
const searchEndpoint = "/search"

type GeocoderClient struct {
	geocoder Geocoder
}

type Place struct {
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	DisplayName string `json:"display_name"`
}

func NewGeoCodeMaps() Geocoder {
	return &GeocodeAPI{baseURL: "https://geocode.maps.co", endpoint: "/search"}
}

func NewGeocoderClient(geocoder Geocoder) *GeocoderClient {
	if geocoder == nil {
		geocoder = &GeocodeAPI{baseURL: defualtBaseURL, endpoint: searchEndpoint}
	}
	return &GeocoderClient{geocoder: geocoder}
}

func (g *GeocodeAPI) fullURL() string {
	return g.baseURL + g.endpoint
}

func (g *GeocoderClient) FindCoordinates(address string) ([]Place, error) {
	queryParams := url.Values{}
	queryParams.Add("q", address)
	// To Do: The addition of the api key to the query params in this way is specific
	// to the geocoder api. So this should happen inside the getPlace method.
	queryParams.Add("api_key", "")
	fullURLtest := fmt.Sprintf("%s?%s", g.geocoder.fullURL(), queryParams.Encode())
	var places []Place
	places, err := g.geocoder.getPlace(fullURLtest)
	return places, err
}

func (g *GeocodeAPI) getPlace(geocodeURL string) ([]Place, error) {
	resp, _ := http.Get(geocodeURL)
	var places []Place

	if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	defer resp.Body.Close()
	return places, nil
}

//func (g *GeocoderClient) FindCoordinates(address string) ([]Place, error) {
//queryParams := url.Values{}
//queryParams.Add("q", address)
//queryParams.Add("api_key", "")
//fullURL := fmt.Sprintf("%s?%s", g.geocoder.fullURL(), queryParams.Encode())
//fmt.Printf("%s", g.geocoder.fullURL())
//resp, _ := http.Get(fullURL)

//var places []Place

//if err := json.NewDecoder(resp.Body).Decode(&places); err != nil {
//fmt.Println("Error decoding JSON:", err)
//}

//defer resp.Body.Close()
//// TO DO: Add proper error response here
//return places, nil
//}
