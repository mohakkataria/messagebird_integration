package main

import (
	"github.com/mohakkataria/messagebird_integration/app"
	"github.com/mohakkataria/messagebird_integration/message_bird"
)

func main() {

	message_bird.InitializeAPIHits()
	app.Start()

}
