package controllers

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "github.com/mohakkataria/messagebird_integration/error"
    "io/ioutil"
    "encoding/json"
    "github.com/mohakkataria/messagebird_integration/util"
    "github.com/mohakkataria/messagebird_integration/models"
    "strings"
    "strconv"
    "github.com/mohakkataria/messagebird_integration/message_bird"
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
        this.Write(w, nil, &error.Error{Code:400, Message:error.BAD_JSON})
        return
    }

    message, error := validateSendMessageAPIInputAndConvertToObject(m)
    if error != nil {
        this.Write(w, nil, error)
        return
    }
    message_bird.QueueMessage(message)

    response := map[string]interface{}{"status":"pending", "message":"message enqueued"}
    this.Write(w, response, nil)
}

func validateSendMessageAPIInputAndConvertToObject(input map[string]interface{}) (*models.Message, *error.Error) {
    message := models.Message{}

    recipients, ok := input["recipient"]
    if !ok {
        return nil, &error.Error{Code:400, Message:error.MISSING_RECIPIENT_INPUT}
    }

    recipientsStringSlice := []string{}
    if util.IsString(recipients) {
        recipientsString := strings.Split(strings.TrimSpace(recipients.(string)), ",")
        for _, recipientString := range recipientsString {
            _, e := strconv.Atoi(strings.TrimSpace(recipientString))
            if e != nil {
                return nil, &error.Error{Code:400, Message:error.BAD_RECIPIENT_INPUT}
            } else {
                recipientsStringSlice = append(recipientsStringSlice, recipientString)
            }
        }
    } else if (util.IsFloat64(recipients)) {
        recipientsStringSlice = append(recipientsStringSlice, strconv.Itoa(int(recipients.(float64))))
    } else {
        return nil, &error.Error{Code:400, Message:error.BAD_RECIPIENT_INPUT}
    }
    message.Recipients = recipientsStringSlice

    originator, ok := input["originator"]
    if !ok {
        return nil, &error.Error{Code:400, Message:error.MISSING_ORIGINATOR_INPUT}
    }
    if util.IsString(originator) {
        originatorString := strings.TrimSpace(originator.(string))
        if util.IsAlphanumeric(originatorString) {
            if len(originatorString) > 11 {
                return nil, &error.Error{Code:400, Message:error.ALPHANUMERIC_LENGTH_ORIGINATOR_ERROR}
            }
            message.Originator = originatorString
        } else {
            return nil, &error.Error{Code:400, Message:error.BAD_ORIGINATOR_INPUT}
        }
    } else if util.IsFloat64(originator) {
        i := int(originator.(float64))
        if (i < 0) {
            return nil, &error.Error{Code:400, Message:error.NUMERIC_ORIGINATOR_ERROR}
        }
        message.Originator = strconv.Itoa(i)
    } else {
        return nil, &error.Error{Code:400, Message:error.BAD_ORIGINATOR_INPUT}
    }

    messageBody, ok := input["message"]
    if !ok {
        return nil, &error.Error{Code:400, Message:error.MISSING_MESSAGE_BODY}
    }
    if util.IsString(messageBody) {
        message.Message = messageBody.(string)
        if util.IsUnicode(message.Message) {
            message.Encoding = models.UNICODE
        } else {
            message.Encoding = models.NORMAL
        }
    } else {
        return nil, &error.Error{Code:400, Message:error.BAD_MESSAGE_INPUT}
    }

    return &message, nil
}
