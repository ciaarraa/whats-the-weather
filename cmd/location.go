/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"whats-the-weather/main/cache"
	"whats-the-weather/main/geocoder"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/zapling/yr.no-golang-client/client"
	"github.com/zapling/yr.no-golang-client/locationforecast"
	"github.com/zapling/yr.no-golang-client/utils"
)

// locationCmd represents the location command
var locationCmd = &cobra.Command{
	Use:   "location [string]",
	Short: "The location you want to see the weather at",
	Long: ` A location to show the weather at. This will find the first possible match and will tell you the
	temperature at that location. Example:
	whats-the-weather location Dubai`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		yrClient := client.NewYrClient(http.DefaultClient, "whats-the-weather-local/0.0 ciaratully0@gmail.com ")
		geoClient := geocoder.NewGeocoderClient(nil)
		cahce_key := args[0]
		cache_db := cache.NewCache("tmp/db", ".cache")

		coords := cache_db.Get(cahce_key)
		var longitude string
		var latitude string

		if coords != "" {
			longLat := strings.Split(coords, ",")
			longitude = longLat[1]
			latitude = longLat[0]
			fmt.Printf("Location: %s\n", args[0])
		} else {
			coordsAndPlaces, _ := geoClient.FindCoordinates(args[0])
			firstMatch := coordsAndPlaces[0]
			longitude = firstMatch.Longitude
			latitude = firstMatch.Latitude
			cache_db.Add([]byte(firstMatch.Latitude+","+firstMatch.Longitude), cahce_key)
			fmt.Printf("Location: %s\n", firstMatch.DisplayName)
		}

		// the following two conversions shoulg be handled in te geocode
		// packaged when the json is being parsed. They should already be
		// float64 objects but for now, we'll convert them here
		floatValueLatitude, err := strconv.ParseFloat(latitude, 64)
		floatValueLongitude, err := strconv.ParseFloat(longitude, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		forecast, _, err := locationforecast.GetCompact(yrClient, floatValueLatitude, floatValueLongitude)
		if err != nil {
			fmt.Print("An error has occured:")
			fmt.Println(err)
			return
		}
		forecastData := forecast.Properties.Timeseries[0].Data
		temperatureNow := forecastData.Instant.Details.AirTemperature
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"", "Next Hour", "Next 6 Hours", "Next 12 Hours"})
		t.AppendRows([]table.Row{{"Weather Summary", *forecastData.Next1Hours.Summary.SymbolCode, *forecastData.Next6Hours.Summary.SymbolCode, *forecastData.Next12Hours.Summary.SymbolCode}})
		t.Render()

		fmt.Printf("temperature at is %v degrees celcius right now\n",  utils.Float64Value(temperatureNow))
	},
}

func init() {
	rootCmd.AddCommand(locationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// locationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// locationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
