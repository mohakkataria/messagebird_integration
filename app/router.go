package app

import (
/*    "net/http"
    "io"*/
    "github.com/julienschmidt/httprouter"
    "github.com/mohakkataria/messagebird_integration/controllers"
)

func (app *App) setRoutes() {

    router := httprouter.New()
    mc := controllers.NewMessageController()

    router.POST("/", mc.SendMessage)

    app.Server.Handler = router
}