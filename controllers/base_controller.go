package controllers

import (
    "encoding/json"
    "net/http"
    "fmt"
)

type BaseController struct {}


func (this BaseController) Write(w http.ResponseWriter, data interface{}, err error) {
    uj, err1 := json.Marshal(data)
    fmt.Println(err1)
    fmt.Println(err)
    // Write content-type, statuscode, payload
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, "%s", uj)
}