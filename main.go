package main

import (
	"github.com/mohakkataria/messagebird_integration/app"
	"github.com/mohakkataria/messagebird_integration/message_bird"
	"github.com/spf13/viper"
	"fmt"
)

func main() {

	viper.SetConfigFile("./config.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded")
	}
	message_bird.InitializeAPIHits()
	app.Start()

}
