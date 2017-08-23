package app

import (
	"net/http"
	"time"
)

var (
	// App is an application instance
	WebApp *App
)

type myHandler struct{}


func init() {
	// create application
    WebApp = NewApp()
}

// App defines application with a new PatternServeMux.
type App struct {
	Server   *http.Server
}

// NewApp returns a new application.
func NewApp() *App {
	app := &App{Server: &http.Server{}}
    app.setRoutes()
	return app
}

// Run application.
func (app *App) Run() {
    addr := ":8080"

    var (
        endRunning = make(chan bool, 1)
    )

	app.Server.ReadTimeout = time.Duration(5) * time.Second
	app.Server.WriteTimeout = time.Duration(5) * time.Second
    app.Server.Handler = &myHandler{}

    go func() {
        app.Server.Addr = addr
        if err := app.Server.ListenAndServe(); err != nil {
            time.Sleep(100 * time.Microsecond)
            endRunning <- true
        }

    }()

	<-endRunning
}

func Start() {
    WebApp.Run()
}

