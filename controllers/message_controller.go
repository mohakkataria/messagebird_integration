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

// MessageController is the controller used for Message related API calls.
type MessageController struct {
    BaseController
}

// NewMessageController returns a pointer to instance of MessageController
func NewMessageController() *MessageController {
    return &MessageController{}
}


// SendMessage accepts the json input, validates the input and Delegates it to message_bird package for enqueuing it
// and sending it upon Rate Limit satisfaction
func (this MessageController) SendMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // Read the body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    m := map[string]interface{}{}
    // Unmarshal it
    err = json.Unmarshal(body, &m)
    if err != nil {
        this.Write(w, nil, &error.Error{Code:400, Message:error.BAD_JSON})
        return
    }

    // validate and get a models.Message object
    message, error := validateSendMessageAPIInputAndConvertToObject(m)
    if error != nil {
        this.Write(w, nil, error)
        return
    }
    // enqueue the models.Message object and delegates the responsibility
    message_bird.QueueMessage(message)
    // since it enqueues, we pass on a pending status to the user
    response := map[string]interface{}{"status":"pending", "message":"message enqueued"}
    this.Write(w, response, nil)
}

// validateSendMessageAPIInputAndConvertToObject takes into account all input related validations
// and formulates a models.Message to be enqueued.
func validateSendMessageAPIInputAndConvertToObject(input map[string]interface{}) (*models.Message, *error.Error) {
    message := models.Message{}

    // check if recipient exists
    recipients, ok := input["recipient"]
    if !ok {
        return nil, &error.Error{Code:400, Message:error.MISSING_RECIPIENT_INPUT}
    }

    recipientsStringSlice := []string{}
    // check if recipients is string
    if util.IsString(recipients) {
        // split the string by comma, and trim the spaces, if any to separate individual recipients
        recipientsString := strings.Split(strings.TrimSpace(recipients.(string)), ",")
        for _, recipientString := range recipientsString {
            _, e := strconv.Atoi(strings.TrimSpace(recipientString))
            if e != nil {
                return nil, &error.Error{Code:400, Message:error.BAD_RECIPIENT_INPUT}
            } else {
                // append recipients to a string slice
                recipientsStringSlice = append(recipientsStringSlice, recipientString)
            }
        }
    } else if (util.IsFloat64(recipients)) {
        // if it is a single number, then convert it to string, and append it to slice
        recipientsStringSlice = append(recipientsStringSlice, strconv.Itoa(int(recipients.(float64))))
    } else {
        return nil, &error.Error{Code:400, Message:error.BAD_RECIPIENT_INPUT}
    }
    message.Recipients = recipientsStringSlice

    // check if originator exists
    originator, ok := input["originator"]
    if !ok {
        return nil, &error.Error{Code:400, Message:error.MISSING_ORIGINATOR_INPUT}
    }
    if util.IsString(originator) {
        // trim spaces if any
        originatorString := strings.TrimSpace(originator.(string))
        if util.IsAlphanumeric(originatorString) {
            // check if length of alphanumeric originator is less than 11
            if len(originatorString) > 11 {
                return nil, &error.Error{Code:400, Message:error.ALPHANUMERIC_LENGTH_ORIGINATOR_ERROR}
            }
            message.Originator = originatorString
        } else {
            return nil, &error.Error{Code:400, Message:error.BAD_ORIGINATOR_INPUT}
        }
    } else if util.IsFloat64(originator) {
        // if it is a single number, check it is not less than 0
        i := int(originator.(float64))
        if (i < 0) {
            return nil, &error.Error{Code:400, Message:error.NUMERIC_ORIGINATOR_ERROR}
        }
        message.Originator = strconv.Itoa(i)
    } else {
        return nil, &error.Error{Code:400, Message:error.BAD_ORIGINATOR_INPUT}
    }

    // check if message body exists
    messageBody, ok := input["message"]
    if !ok {
        return nil, &error.Error{Code:400, Message:error.MISSING_MESSAGE_BODY}
    }
    if util.IsString(messageBody) {
        // check the encoding of the message body
        message.MessageBody = messageBody.(string)
        if util.IsUnicode(message.MessageBody) {
            message.Encoding = models.UNICODE
        } else {
            message.Encoding = models.NORMAL
        }
    } else {
        return nil, &error.Error{Code:400, Message:error.BAD_MESSAGE_INPUT}
    }

    return &message, nil
}
