package geocoder

type Geocoder interface {
	getPlace(address string) ([]Place, error)
	fullURL() string
}

type Place struct {
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	DisplayName string `json:"display_name"`
}
