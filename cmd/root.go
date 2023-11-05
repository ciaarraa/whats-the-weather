package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"whats-the-weather/main/geocoder"

	"github.com/spf13/cobra"
	"github.com/zapling/yr.no-golang-client/client"
	"github.com/zapling/yr.no-golang-client/locationforecast"
	"github.com/zapling/yr.no-golang-client/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whats-the-weather [string]",
	Short: "What is the weather?",
	Long: `A cli tool which tells you the weather in a location. The defaut location
	is Big Tesco in Southwark, London.`,
	Args: cobra.ExactArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		yrClient := client.NewYrClient(http.DefaultClient, "Yr.no golang example")
		//var defaultLongitude, defaultLatitude float64
		//defaultLongitude = 51.501820
		//defaultLatitude = -0.097920

		geoClient := geocoder.NewGeocoderClient()
		coordsAndPlaces, _ := geoClient.FindCoordinates(args[0])
		firstMatch := coordsAndPlaces[0]
		// the following two conversions shoulg be handled in te geocode
		// packaged when the json is being parsed. They should already be
		// float64 objects but for now, we'll convert them here
		floatValueLatitude, err := strconv.ParseFloat(firstMatch.Latitude, 64)
		floatValueLongitude, err := strconv.ParseFloat(firstMatch.Longitude, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// here we want to implemnt a cache. If there is a cached forecast for the longitude and latitude,
		// find that instead of calling the api again.
		forecast, _, err := locationforecast.GetCompact(yrClient, floatValueLatitude, floatValueLongitude)
		if err != nil {
			fmt.Print("An error has occured:")
			fmt.Println(err)
			return
		}

		temperature := forecast.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
		fmt.Printf("temperature at %s is : %v degrees celcius\n", firstMatch.DisplayName, utils.Float64Value(temperature))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
