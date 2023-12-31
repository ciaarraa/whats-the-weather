package geocoder

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlace(t *testing.T) {
	type server_config struct {
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

	createServerMock := func(server_config_args server_config, serverResponseBody []ServerResponse) *httptest.Server {
		s := httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(server_config_args.httpGetStatus)
					jsonResponse, _ := json.Marshal(serverResponseBody)
					w.Write(jsonResponse)
				}),
		)
		return s
	}

	tests := []struct {
		name               string
		args               server_config
		serverJSONResponse []ServerResponse
		want               []Place
		wantErr            error
	}{
		{
			name:               "Returns slice of Places when 200 response",
			args:               server_config{httpGetStatus: 200},
			serverJSONResponse: []ServerResponse{{Lat: "100", Lon: "100", DisplayName: "Somewhere", PlaceID: 123}, {Lat: "200", Lon: "100", DisplayName: "Everest", PlaceID: 789}},
			want:               []Place{{Latitude: "100", Longitude: "100", DisplayName: "Somewhere"}, {Latitude: "200", Longitude: "100", DisplayName: "Everest"}},
		},
		{
			name:               "Returns empty Place slice if the geocode returns no results ",
			args:               server_config{httpGetStatus: 200},
			serverJSONResponse: []ServerResponse{},
			want:               []Place{},
		},
		{
			name:               "Returns empty Place slice if geocode returns a 500",
			args:               server_config{httpGetStatus: 500},
			serverJSONResponse: []ServerResponse{},
			want:               []Place{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serverMock := createServerMock(tt.args, tt.serverJSONResponse)
			geo := GeocodeAPI{baseURL: serverMock.URL, endpoint: "/search", apiKey: "testKey"}
			got, err := geo.getPlace("test_location")
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
			serverMock.Close()
		})
	}
}
