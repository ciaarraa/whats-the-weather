/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/zapling/yr.no-golang-client/client"
	"github.com/zapling/yr.no-golang-client/locationforecast"
	"github.com/zapling/yr.no-golang-client/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whats-the-weather",
	Short: "What is the weather?",
	Long: `A cli tool which tells you the weather in a location. The defaut location
	is Big Tesco in Southwark, London. For other locations ask for help`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		yrClient := client.NewYrClient(http.DefaultClient, "Yr.no golang example")

		forecast, _, err := locationforecast.GetCompact(yrClient, 51.501820, -0.097920)

		if err != nil {
			fmt.Println(err)
			return
		}

		temperature := forecast.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
		fmt.Printf("temperature: %v\n", utils.Float64Value(temperature))
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
