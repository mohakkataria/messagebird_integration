package main

import (
	"fmt"
	"github.com/mohakkataria/messagebird_integration/app"
	"github.com/mohakkataria/messagebird_integration/message_bird"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile("./config.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded")
	}
	messageBird.InitializeAPIHits()
	app.Start()

}
