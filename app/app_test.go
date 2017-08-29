package app

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func TestStart(t *testing.T) {
	go Start()
	addressInfo := viper.GetStringMapString("addressInfo")
	addr := addressInfo["host"] + ":" + addressInfo["port"]
	res, _ := http.Get("http://" + addr + "/")
	fmt.Println(addr + "/")
	if res.StatusCode != 405 {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", 405, res.StatusCode)
	}

}

func init() {
	viper.SetConfigFile("./../config.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded")
	}
}
