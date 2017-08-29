package main

import (
	"fmt"
	"github.com/mohakkataria/messagebird_integration/app"
	"github.com/mohakkataria/messagebird_integration/messageBird"
	"github.com/spf13/viper"
)

func main() {
	app.Start()
}

func init() {
	viper.SetConfigFile("./config.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded")
	}
	messageBird.Initialize()
	messageBird.StartChannelConsumer()
}
