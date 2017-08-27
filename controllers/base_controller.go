package controllers

import (
    "encoding/json"
    "net/http"
    "fmt"
    "github.com/mohakkataria/messagebird_integration/error"
)

type BaseController struct {}


func (this BaseController) Write(w http.ResponseWriter, data interface{}, err *error.Error) {
    w.Header().Set("Content-Type", "application/json")
    if err != nil {
        w.WriteHeader(err.Code)
        errMap := map[string]string{"status" : "failed", "error" : err.Message}
        errString, _ := json.Marshal(errMap)
        fmt.Fprintf(w, "%s", errString)
        fmt.Println(errMap)
        return
    }

    uj, err1 := json.Marshal(data)
    fmt.Println(err1)
    fmt.Println(err)
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}