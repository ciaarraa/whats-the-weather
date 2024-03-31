package forecastcache

import (
	"encoding/json"
	"fmt"
	"time"
	"whats-the-weather/main/cache"

	"github.com/zapling/yr.no-golang-client/locationforecast"
)

type funcTimeType func() time.Time

type ForecastCache struct {
	Cache *cache.Cache
	TimeNow funcTimeType
}

type cachedForecast struct{
	Forecast locationforecast.GeoJson
	ExpiresAt time.Time
}

func (forecastCache *ForecastCache) Add(forecast locationforecast.GeoJson, until time.Time, key string) {
	cachedObj := cachedForecast{Forecast: forecast, ExpiresAt: until}
	jsonCachedObj, _ := json.Marshal(cachedObj)
	var cachedObjNew locationforecast.GeoJson
	json.Unmarshal(jsonCachedObj, &cachedObjNew)

	forecastCache.Cache.Add(jsonCachedObj, key)
}

func (forecastCache *ForecastCache) Get(key string) (locationforecast.GeoJson, error){
	cachedForecastObj := forecastCache.Cache.Get(key)
	var cachedObjNew cachedForecast
	json.Unmarshal([]byte(cachedForecastObj), &cachedObjNew)
	fmt.Print(cachedObjNew.ExpiresAt)
	checkExpiryAgainst := forecastCache.TimeNow()
	if cachedForecastObj == "" {
		return locationforecast.GeoJson{} , nil
	}
	if cachedObjNew.ExpiresAt.Sub(checkExpiryAgainst) > 0 {
		return cachedObjNew.Forecast, nil
	}else{
		return locationforecast.GeoJson{} , nil
	}
}


