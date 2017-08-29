package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/mohakkataria/messagebird_integration/error"
	"github.com/mohakkataria/messagebird_integration/messageBird"
	"github.com/mohakkataria/messagebird_integration/models"
	"github.com/mohakkataria/messagebird_integration/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// MessageController is the controller used for Message related API calls.
type MessageController struct {
}

// NewMessageController returns a pointer to instance of MessageController
func NewMessageController() *MessageController {
	return &MessageController{}
}

// SendMessage accepts the json input, validates the input and Delegates it to message_bird package for enqueuing it
// and sending it upon Rate Limit satisfaction
func (mc MessageController) SendMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	m := map[string]interface{}{}
	// Unmarshal it
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Println(err)
		Write(w, nil, &error.Error{Code: 400, Message: error.BadJSON})
		return
	}

	// validate and get a models.Message object
	message, error := validateSendMessageAPIInputAndConvertToObject(m)
	if error != nil {
		Write(w, nil, error)
		return
	}
	// enqueue the models.Message object and delegates the responsibility
	messageBird.QueueMessage(message)
	// since it enqueues, we pass on a pending status to the user
	response := map[string]interface{}{"status": "pending", "message": "message enqueued"}
	Write(w, response, nil)
}

// validateSendMessageAPIInputAndConvertToObject takes into account all input related validations
// and formulates a models.Message to be enqueued.
func validateSendMessageAPIInputAndConvertToObject(input map[string]interface{}) (*models.Message, *error.Error) {
	message := models.Message{}

	// check if recipient exists
	recipients, ok := input["recipient"]
	if !ok {
		return nil, &error.Error{Code: 400, Message: error.MissingRecipientInput}
	}

	recipientsStringSlice := []string{}
	// check if recipients is string
	if util.IsString(recipients) {
		// split the string by comma, and trim the spaces, if any to separate individual recipients
		recipientsString := strings.Split(strings.TrimSpace(recipients.(string)), ",")
		for _, recipientString := range recipientsString {
			_, e := strconv.Atoi(strings.TrimSpace(recipientString))
			if e != nil {
				return nil, &error.Error{Code: 400, Message: error.BadRecipientInput}
			}
			// append recipients to a string slice
			recipientsStringSlice = append(recipientsStringSlice, recipientString)
		}
	} else if util.IsFloat64(recipients) {
		// if it is a single number, then convert it to string, and append it to slice
		recipientsStringSlice = append(recipientsStringSlice, strconv.Itoa(int(recipients.(float64))))
	} /*else {
		return nil, &error.Error{Code: 400, Message: error.BAD_RECIPIENT_INPUT}
	}*/
	message.Recipients = recipientsStringSlice

	// check if originator exists
	originator, ok := input["originator"]
	if !ok {
		return nil, &error.Error{Code: 400, Message: error.MissingOriginatorInput}
	}
	if util.IsString(originator) {
		// trim spaces if any
		originatorString := strings.TrimSpace(originator.(string))
		if util.IsAlphanumeric(originatorString) {
			// check if length of alphanumeric originator is less than 11
			if len(originatorString) > 11 {
				return nil, &error.Error{Code: 400, Message: error.AlphanumericLengthOriginatorError}
			}
			message.Originator = originatorString
		} else {
			return nil, &error.Error{Code: 400, Message: error.BadOriginatorInput}
		}
	} else if util.IsFloat64(originator) {
		// if it is a single number, check it is not less than 0
		i := int(originator.(float64))
		if i < 0 {
			return nil, &error.Error{Code: 400, Message: error.NumericOriginatorInput}
		}
		message.Originator = strconv.Itoa(i)
	} /*else {
		return nil, &error.Error{Code: 400, Message: error.BAD_ORIGINATOR_INPUT}
	}*/

	// check if message body exists
	messageBody, ok := input["message"]
	if !ok {
		return nil, &error.Error{Code: 400, Message: error.MissingMessageBody}
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
		return nil, &error.Error{Code: 400, Message: error.BadMessageInput}
	}

	return &message, nil
}
