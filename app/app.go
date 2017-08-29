package app

import (
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	// WebApp is an HTTP Application instance
	WebApp *App
)

func init() {
	WebApp = NewApp()
}

// App defines http server application
type App struct {
	Server *http.Server
}

// NewApp returns a new application instance
func NewApp() *App {
	app := &App{Server: &http.Server{}}
	app.setRoutes()
	return app
}

// Run function, runs the HTTP Server
func (app *App) run() {
	addressInfo := viper.GetStringMapString("addressInfo")
	addr := addressInfo["host"] + ":" + addressInfo["port"]

	var (
		endRunning = make(chan bool, 1)
	)

	app.Server.ReadTimeout = time.Duration(5) * time.Second
	app.Server.WriteTimeout = time.Duration(5) * time.Second

	go func() {
		app.Server.Addr = addr
		if err := app.Server.ListenAndServe(); err != nil {
			time.Sleep(100 * time.Microsecond)
			endRunning <- true
		}

	}()

	<-endRunning
}

// Start function is exported for other packages to call. It is called in main
func Start() {
	WebApp.run()
}
