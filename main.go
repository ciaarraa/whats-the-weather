/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"whats-the-weather/main/cmd"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()

}
