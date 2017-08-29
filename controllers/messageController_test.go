package controllers

import (
    "testing"
    "encoding/json"
    "github.com/mohakkataria/messagebird_integration/error"
    "net/http/httptest"
    "strings"
    "github.com/mohakkataria/messagebird_integration/message_bird"
    "reflect"
    "github.com/spf13/viper"
    "fmt"
    "github.com/mohakkataria/messagebird_integration/messageBird"
)

func Test_validateSendMessageAPIInputAndConvertToObject(t *testing.T) {
    msgInput := []byte(`{"recipient":123,"originator":123,"message" : "test"}`)
    msg := map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e := validateSendMessageAPIInputAndConvertToObject(msg)
    if e != nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", nil, e.Message)
    }

    msgInput = []byte(`{"recipient":123,"originator":"asdasdasdasd","message" : "test"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.alphanumericLengthOriginatorError, nil)
    }

    msgInput = []byte(`{"recipient":"123","originator":"123","message" : "test"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e != nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", nil, e.Message)
    }

    msgInput = []byte(`{"recipient":"123,123","originator":"asdasdasdasd","message" : "test"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.AlphanumericLengthOriginatorError, nil)
    }

    msgInput = []byte(`{"recipient":"123,123","originator":"!!!","message" : "test"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.BadOriginatorInput, nil)
    }

    msgInput = []byte(`{"recipient":"123,123","message" : "test"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.MissingOriginatorInput, nil)
    }

    msgInput = []byte(`{"originator":"!!!","message" : "test"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.MissingRecipientInput, nil)
    }

    msgInput = []byte(`{"recipient":"123,123","originator":"!!!"}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.MissingMessageBody, nil)
    }

    msgInput = []byte(`{"recipient":"123,123","originator":"123","message" : 123}`)
    msg = map[string]interface{}{}
    json.Unmarshal(msgInput, &msg)
    _, e = validateSendMessageAPIInputAndConvertToObject(msg)
    if e == nil {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", error.BadMessageInput, nil)
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
        t.Errorf("Test failed, expected: '%s', got:  '%s'", 200, w.Code)
    }

    payload = `{"recipient":123,"originator":"!!","message" : "test"}`
    req = httptest.NewRequest("POST", "/", strings.NewReader(payload))
    w = httptest.NewRecorder()
    mc.SendMessage(w, req, nil)
    if w.Code != 400 {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", 400, w.Code)
    }
}

func TestMessageController_SendMessage2(t *testing.T) {
    mc := MessageController{}
    req := httptest.NewRequest("POST", "/", nil)
    w := httptest.NewRecorder()
    mc.SendMessage(w, req, nil)
    if w.Code != 400 {
        t.Errorf("Test failed, expected: '%s', got:  '%s'", 400, w.Code)
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