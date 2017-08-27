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
        this.Write(w, nil, &error.Error{Code:400, Message:"Bad Request Json"})
        return
    }

    _, error := validateSendMessageAPIInputAndConvertToObject(m)
    if error != nil {
        this.Write(w, nil, error)
        return
    }

    this.Write(w, m, nil)
}

func validateSendMessageAPIInputAndConvertToObject(input map[string]interface{}) (*models.Message, *error.Error) {
    message := models.Message{}

    recipients, ok := input["recipient"]
    if !ok {
        return nil, &error.Error{Code:400, Message:"Missing Recipient"}
    }

    recipientsInteger := []int64{}
    if util.IsString(recipients) {
        recipientsString := strings.Split(strings.TrimSpace(recipients.(string)), ",")
        for _, recipientString := range recipientsString {
            i, e := strconv.Atoi(strings.TrimSpace(recipientString))
            if e != nil {
                return nil, &error.Error{Code:400, Message:"Bad Recipient Input"}
            } else {
                recipientsInteger = append(recipientsInteger, int64(i))
            }
        }
    } else if (util.IsFloat64(recipients)) {
        recipientsInteger = append(recipientsInteger, int64(recipients.(float64)))
    } else {
        return nil, &error.Error{Code:400, Message:"Bad Recipient Input"}
    }
    message.Recipients = recipientsInteger

    originator, ok := input["originator"]
    if !ok {
        return nil, &error.Error{Code:400, Message:"Missing Originator"}
    }
    if util.IsString(originator) {
        originatorString := strings.TrimSpace(originator.(string))
        if util.IsAlphanumeric(originatorString) {
            if len(originatorString) > 11 {
                return nil, &error.Error{Code:400, Message:"Alphanumeric Originator should be less than 11 characters"}
            }
            message.Originator = originatorString
        }
        return nil, &error.Error{Code:400, Message:"Bad Originator input"}

    } else if util.IsFloat64(originator) {
        i := int(originator.(float64))
        if (i < 0) {
            return nil, &error.Error{Code:400, Message:"Bad Originator Input as it cannot be a number less than 0"}
        }
        message.Originator = strconv.Itoa(i)
    } else {
        return nil, &error.Error{Code:400, Message:"Bad Originator input"}
    }

    messageBody, ok := input["message"]
    if !ok {
        return nil, &error.Error{Code:400, Message:"Missing Message Body"}
    }
    if util.IsString(messageBody) {
        message.Message = messageBody.(string)
        if util.IsUnicode(message.Message) {
            message.Encoding = models.UNICODE
        } else {
            message.Encoding = models.NORMAL
        }
    } else {
        return nil, &error.Error{Code:400, Message:"Bad Message Body"}
    }

    return &message, nil
}
