package geocoder

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCoordinates(t *testing.T) {
	type args struct {
		url           string
		httpGetStatus int
	}

	// the json keys returned from the geocoder api
	type ServerResponse struct {
		Lat         string   `json:"lat"`
		Lon         string   `json:"lon"`
		DisplayName string   `json:"display_name"`
		Class       string   `json:"class"`
		Type        string   `json:"type"`
		Importance  float64  `json:"importance"`
		PlaceID     int      `json:"place_id"`
		License     string   `json:"licence"`
		PoweredBy   string   `json:"powered_by"`
		OSMType     string   `json:"osm_type"`
		OSMID       int      `json:"osm_id"`
		BoundingBox []string `json:"boundingbox"`
	}

	createServerMock := func(args2 args, serverResponseBody []ServerResponse) *httptest.Server {
		s := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(args2.httpGetStatus)
					jsonResponse, _ := json.Marshal(serverResponseBody)
					w.Write(jsonResponse)
				}),
		)
		return s
	}

	tests := []struct {
		name               string
		args               args
		serverJSONResponse []ServerResponse // TO DO: refactor this to move it into args
		want               []Place
		wantErr            error
	}{
		{
			name:               "Returns slice of Places when 200 response",
			args:               args{httpGetStatus: 200},
			serverJSONResponse: []ServerResponse{{Lat: "100", Lon: "100", DisplayName: "Somewhere", PlaceID: 123}, {Lat: "200", Lon: "100", DisplayName: "Everest", PlaceID: 789}},
			want:               []Place{{Latitude: "100", Longitude: "100", DisplayName: "Somewhere"}, {Latitude: "200", Longitude: "100", DisplayName: "Everest"}},
		},
		{
			name:               "Returns empty Place slice if the geocode returns no results ",
			args:               args{httpGetStatus: 200},
			serverJSONResponse: []ServerResponse{},
			want:               []Place{},
		},
		{
			name:               "Returns empty Place slice if geocode returns a 500",
			args:               args{httpGetStatus: 500},
			serverJSONResponse: []ServerResponse{},
			want:               []Place{},
		},
	}

	geo := GeocodeAPI{baseURL: "testBase.com", endpoint: "/search"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverMock := createServerMock(tt.args, tt.serverJSONResponse)
			tt.args.url = serverMock.URL
			got, err := geo.getPlace(tt.args.url)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
			serverMock.Close()
		})
	}
}

type MockGeocodeAPI struct {
	places []Place
	err    error
}

func (m MockGeocodeAPI) getPlace(address string) ([]Place, error) {
	return m.places, m.err
}
func (m MockGeocodeAPI) fullURL() string { return "testing" }

func TestFindLocation(t *testing.T) {
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
			mockGeocodeAPI := MockGeocodeAPI{places: tt.args.mockResponse}
			geocoderClient := &GeocoderClient{geocoder: mockGeocodeAPI}
			location, _ := geocoderClient.FindLocation("test address")
			assert.Equal(t, tt.want, location)

		})
	}
}
