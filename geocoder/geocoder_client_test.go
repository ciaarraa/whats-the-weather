package geocoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockGeocode struct {
	places []Place
	err    error
}

func (m MockGeocode) getPlace(address string) ([]Place, error) {
	return m.places, m.err
}
func (m MockGeocode) fullURL() string { return "testing" }

func TestFindCoordinates(t *testing.T) {
	type args struct {
		mockResponse []Place
	}
	tests := []struct {
		name    string
		args    args
		want    []Place
		wantErr error
	}{
		{
			name: "It delegates to the geocoder api and returns the response and error",
			args: args{mockResponse: []Place{{Latitude: "100", Longitude: "100", DisplayName: "Somewhere"}}},
			want: []Place{{Latitude: "100", Longitude: "100", DisplayName: "Somewhere"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGeocodeAPI := MockGeocode{places: tt.args.mockResponse}
			geocoderClient := &GeocoderClient{geocoder: mockGeocodeAPI}
			location, _ := geocoderClient.FindCoordinates("test address")
			assert.Equal(t, tt.want, location)

		})
	}
}
