package app

import (
    "net/http"
    "io"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if h, ok := mux[r.URL.String()]; ok {
        h(w, r)
        return
    }

    io.WriteString(w, "My server: " + r.URL.String())
}

func hello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello world!")
}

func (*App) setRoutes() {
    mux = make(map[string]func(http.ResponseWriter, *http.Request))
    mux["/"] = hello
}