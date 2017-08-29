package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/mohakkataria/messagebird_integration/error"
	"github.com/mohakkataria/messagebird_integration/messageBird"
	"github.com/mohakkataria/messagebird_integration/message_bird"
	"github.com/spf13/viper"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_validateSendMessageAPIInputAndConvertToObject(t *testing.T) {
	msgInput := []byte(`{"recipient":123,"originator":123,"message" : "test"}`)
	msg := map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e := validateSendMessageAPIInputAndConvertToObject(msg)
	if e != nil {
		t.Errorf("Test failed, expected: nil, got:  '%s'", e.Message)
	}

	msgInput = []byte(`{"recipient":123,"originator":"asdasdasdasd","message" : "test"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got: nil", error.AlphanumericLengthOriginatorError)
	}

	msgInput = []byte(`{"recipient":"123","originator":"123","message" : "test"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e != nil {
		t.Errorf("Test failed, expected: nil, got:  '%s'", e.Message)
	}

	msgInput = []byte(`{"recipient":"123,123","originator":"asdasdasdasd","message" : "test"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got:  nil", error.AlphanumericLengthOriginatorError)
	}

	msgInput = []byte(`{"recipient":"123,123","originator":"!!!","message" : "test"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got:  nil", error.BadOriginatorInput)
	}

	msgInput = []byte(`{"recipient":"123,123","message" : "test"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got:  nil", error.MissingOriginatorInput)
	}

	msgInput = []byte(`{"originator":"!!!","message" : "test"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got:  nil", error.MissingRecipientInput)
	}

	msgInput = []byte(`{"recipient":"123,123","originator":"!!!"}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got:  nil", error.MissingMessageBody)
	}

	msgInput = []byte(`{"recipient":"123,123","originator":"123","message" : 123}`)
	msg = map[string]interface{}{}
	json.Unmarshal(msgInput, &msg)
	_, e = validateSendMessageAPIInputAndConvertToObject(msg)
	if e == nil {
		t.Errorf("Test failed, expected: '%s', got: nil", error.BadMessageInput)
	}
}

func TestMessageController_SendMessage(t *testing.T) {
	messageBird.InitializeAPIHits()
	mc := MessageController{}
	payload := `{"recipient":123,"originator":123,"message" : "test"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(payload))
	w := httptest.NewRecorder()
	mc.SendMessage(w, req, nil)
	if w.Code != 200 {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", 200, w.Code)
	}

	payload = `{"recipient":123,"originator":"!!","message" : "test"}`
	req = httptest.NewRequest("POST", "/", strings.NewReader(payload))
	w = httptest.NewRecorder()
	mc.SendMessage(w, req, nil)
	if w.Code != 400 {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", 400, w.Code)
	}
}

func TestMessageController_SendMessage2(t *testing.T) {
	mc := MessageController{}
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	mc.SendMessage(w, req, nil)
	if w.Code != 400 {
		t.Errorf("Test failed, expected: '%d', got:  '%d'", 400, w.Code)
	}
}

func TestNewMessageController(t *testing.T) {
	mc := &MessageController{}
	mc1 := NewMessageController()
	if !reflect.DeepEqual(mc1, mc) {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", mc1, mc)
	}
}

func init() {
	viper.SetConfigFile("./../config.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No configuration file loaded")
	}
}
