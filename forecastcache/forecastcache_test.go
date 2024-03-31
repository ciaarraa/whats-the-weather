package locationcache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"whats-the-weather/main/cache"

	"github.com/stretchr/testify/assert"
	"github.com/zapling/yr.no-golang-client/locationforecast"
)

func futureTime() time.Time{
	return time.Date(2998, time.January, 5, 0, 0, 0, 0, time.UTC)
}

func pastTime() time.Time{
	return time.Date(1998, time.January, 5, 0, 0, 0, 0, time.UTC)
}


func TestLocationCache(t *testing.T) {
	var forecast locationforecast.GeoJson
	tempDir := t.TempDir()
	dbPath := tempDir + "/db"
	cacheFolder := tempDir + "/.cache"
	testCache := cache.NewCache(dbPath, cacheFolder)

	tests := []struct {
		name string
		expiryTime time.Time
		mockForecastResponseFile string
		want string
		timeFunc time.Time
		cacheKey string

		}{
			{
				name: "Expired Cache",
				expiryTime: time.Date(1998, time.January, 5, 0, 0, 0, 0, time.UTC),
				mockForecastResponseFile: "../example_response.txt",
				timeFunc: futureTime(),
				cacheKey: "cacheKey",
			},
			{
				name: "Valid Cache",
				expiryTime: time.Date(2998, time.January, 5, 0, 0, 0, 0, time.UTC),
				mockForecastResponseFile: "../example_response.txt",
				timeFunc: pastTime(),
				cacheKey: "secondCacheKey",

			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
			getTimeFunc := func() time.Time {
					return tt.timeFunc
				}
			forecastCache := ForecastCache{cache: testCache, timeNow: getTimeFunc }
			expiryTime := tt.expiryTime
			forecastResponse, err := os.Open(tt.mockForecastResponseFile)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer forecastResponse.Close()

			decoder := json.NewDecoder(forecastResponse)
			if err := decoder.Decode(&forecast); err != nil {
				log.Fatalf("Error decoding JSON: %v", err)
			}
			fmt.Print(expiryTime)
			forecastCache.Add(forecast, expiryTime, tt.cacheKey)

			returnedFromCache, _ := forecastCache.Get(tt.cacheKey)
			if forecastTimeseries := returnedFromCache.Properties.Timeseries; len(forecastTimeseries) > 0 {
				assert.Equal(t, forecast, returnedFromCache )
			} else {
				fmt.Println("In here")
			}
		})
	}}