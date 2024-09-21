package cmd

import (
	"os"
	"github.com/spf13/cobra"
)
var debug bool
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "whats-the-weather [string]",
	Short: "What is the weather?",
	Long: `A cli tool which tells you the weather in a location.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	PersistentPreRun: toggleDebug,
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
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false,  "Show debug logs")
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
