package geocoder

import (
	"os"
)

const defualtBaseURL = "https://geocode.maps.co"
const searchEndpoint = "/search"

type GeocoderClient struct {
	geocoder Geocoder
}

func NewGeocoderClient(geocoder Geocoder) *GeocoderClient {
	if geocoder == nil {
		geocoder = &GeocodeAPI{baseURL: defualtBaseURL, endpoint: searchEndpoint, apiKey: os.Getenv("GEOCODE_API_KEY")}
	}
	return &GeocoderClient{geocoder: geocoder}
}

func (g *GeocoderClient) FindCoordinates(address string) ([]Place, error) {
	var places []Place
	places, err := g.geocoder.getPlace(address)
	return places, err
}
