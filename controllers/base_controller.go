// Package controllers contains the controllers used to serve API requests
package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/mohakkataria/messagebird_integration/error"
	"net/http"
)

// BaseController is an empty struct used to declare the Write functionality to be used by all
// other controllers. This is a very naive implementation of controller for the purpose of this exercise
type BaseController struct{}

// Write function, write the error json with appropriate code or the data as passed on by the controller calling it
func (this BaseController) Write(w http.ResponseWriter, data interface{}, err *error.Error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(err.Code)
		errMap := map[string]string{"status": "failed", "error": err.Message}
		errString, _ := json.Marshal(errMap)
		fmt.Fprintf(w, "%s", errString)
		fmt.Println(errMap)
		return
	}

	uj, _ := json.Marshal(data)
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
