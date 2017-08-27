package controllers

import (
    "net/http"
    "github.com/julienschmidt/httprouter"

    "io/ioutil"
    "encoding/json"
)

type MessageController struct {
    BaseController
}

func NewMessageController() *MessageController {
    return &MessageController{}
}

func (this MessageController) SendMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    m := map[string]interface{}{}
    err = json.Unmarshal(body, &m)
    if err != nil {
        this.Write(w, nil, err)
        return
    }

    this.Write(w, m, nil)
}

