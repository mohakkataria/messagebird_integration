// Package message_bird contains the functionality to interact with
// MessageBird SDK
package message_bird

import (
    "github.com/messagebird/go-rest-api"
    "github.com/mohakkataria/messagebird_integration/models"
    "github.com/mohakkataria/messagebird_integration/util"
    "time"
    "fmt"
)

// sendSingleMessageRequests channel is used hold the requests to be sent
// via Messagebird API as we are Rate Limitting the output to API
var sendSingleMessageRequests chan models.SplitMessage

// messageBirdClient Message Bird Client
var messageBirdClient *messagebird.Client

// Constants for various message length restrictions and API
const (
    MAX_MSG_SIZE_PLAIN = 1377
    MAX_SINGLE_MSG_SEGMENT_SIZE_PLAIN = 160
    MAX_MULTIPART_MSG_SEGMENT_SIZE_PLAIN = 153
    MAX_MSG_SIZE_UNICODE = 603
    MAX_SINGLE_MSG_SEGMENT_SIZE_UNICODE = 70
    MAX_MULTIPART_MSG_SEGMENT_SIZE_UNICODE = 67
    API_KEY = "apwitDmPqmha2r4OdhAhFDvAwCOPY"
    API_RATE_LIMIT = 1
)

// QueueMessage is exposed to outside the package for the input Message object to be enqueued
// for sending via Message Bird SDK after processing it for Message Body checks and adding
// appropriate UDH if required
func QueueMessage(message * models.Message) error {
    messageBodyLength := message.GetMessagebodyLength()

    // Get splitMessage from Message input
    splitMessage := message.GetSplitMessageWithOutBodyFromMessage()
    // Here we  are making sure that the number of characters are as per specification.
    // Segmenting logic is also handled by dividing the message into chunks if required.
    if (!message.IsEncodingNormal()) {
        if (messageBodyLength < MAX_SINGLE_MSG_SEGMENT_SIZE_UNICODE) {
            splitMessage.MessageBodyChunk = message.MessageBody
            queueSplitMessage(splitMessage)
        } else {
            // Truncate message if length is more than the allowed length
            if (messageBodyLength > MAX_MSG_SIZE_UNICODE) {
                message.MessageBody = string([]rune(message.MessageBody)[:MAX_MSG_SIZE_UNICODE])
            }

            // Get chunks
            messages := getMessageChunks(MAX_MSG_SIZE_UNICODE, MAX_MULTIPART_MSG_SEGMENT_SIZE_UNICODE, true, message.MessageBody)
            enqueueConcatenatedSMS(splitMessage, messages)
        }
    } else {
        if (messageBodyLength < MAX_SINGLE_MSG_SEGMENT_SIZE_PLAIN) {
            splitMessage.MessageBodyChunk = message.MessageBody
            fmt.Println(splitMessage)
            queueSplitMessage(splitMessage)
        } else {
            // if message exceeds maximum size limit for plain messages, we only send the maximum allowed chars.
            if (messageBodyLength > MAX_MSG_SIZE_PLAIN) {
                message.MessageBody = message.MessageBody[:MAX_MULTIPART_MSG_SEGMENT_SIZE_PLAIN]
            }

            // Get chunks
            messages := getMessageChunks(MAX_MSG_SIZE_PLAIN, MAX_MULTIPART_MSG_SEGMENT_SIZE_PLAIN, false, message.MessageBody)
            enqueueConcatenatedSMS(splitMessage, messages)
        }
    }
    return nil;
}

// enqueueConcatenatedSMS enqueues UDH marked Concatenated SMSes for sending
func enqueueConcatenatedSMS(splitMessage models.SplitMessage, chunks []string) {
    segments := 0
    numberOfSegments := len(chunks)
    // Generate random 2 digit hex reference number
    referenceNumber := util.GetRandom2DigitHex()
    for _, body := range chunks {
        segments++;
        splitMessage.MessageBodyChunk = body
        /*
            UDH format
            Field 1 (1 octet): Length of User Data Header, in this case 05.
            Field 2 (1 octet): Information Element Identifier, equal to 00
            Field 3 (1 octet): Length of the header, excluding the first two fields; equal to 03
            Field 4 (1 octet): 00-FF, CSMS reference number, must be same for all the SMS parts in the CSMS
            Field 5 (1 octet): 00-FF, total number of parts.
            Field 6 (1 octet): 00-FF, this part's number in the sequence. The value shall start at 1.
        */
        // Mark UDH for the message
        splitMessage.UDH = "050003" + referenceNumber + util.ConvertIntToHex(numberOfSegments) + util.ConvertIntToHex(segments);
        queueSplitMessage(splitMessage)
    }
}

// getMessageChunks gets a chunk of message body depending upon the encoding of the message
func getMessageChunks(maxLimit int, chunkSize int, unicode bool, messageBody string) []string {
    messages := []string{}
    i := 0
    // divide body into segments of specified size.
    for (i < maxLimit) {
        if unicode {
            messages = append(messages, string([]rune(messageBody)[i:i+chunkSize]))
        } else {
            messages = append(messages, messageBody[i:i+chunkSize])
        }
        i += chunkSize
    }
    return messages
}

// queueSplitMessage enqueues a splitMessage to sendSingleMessageRequests channel
func queueSplitMessage(splitMessage models.SplitMessage) {
    sendSingleMessageRequests <- splitMessage
}

// sendSingleMessage sends a message vis MessageBird SDK
func sendSingleMessage(message models.SplitMessage) {
    // convert the SplitMessage to MessageParams and then retrieve the remaining properties and make a call to
    // message bird API.
    messageParams := &messagebird.MessageParams{TypeDetails:map[string]interface{}{"udh" : message.UDH}, DataCoding:message.DataCoding}
    mbMessage, err := messageBirdClient.NewMessage(message.Originator, message.Recipients, message.MessageBodyChunk, messageParams)
    if err != nil {
        // messagebird.ErrResponse means custom JSON errors.
        if err == messagebird.ErrResponse {
            for _, mbError := range mbMessage.Errors {
                fmt.Printf("Error: %#v\n", mbError)
            }
        }
    }

}


// InitializeAPIHits initializes the required channel for send requests
// and the message bird client to be used. Also we initialize a ticker of 1 second
// to implement the Rate Limit as Mentioned. As discussed, we are not using a Redis Based approach
// and this will only be limited to this app instance
func InitializeAPIHits() {
    sendSingleMessageRequests = make(chan models.SplitMessage);
    messageBirdClient = messagebird.New(API_KEY)
    rate := time.Second / API_RATE_LIMIT
    throttle := time.Tick(rate)
    fmt.Println("+1")
    go func() {
        for  {
            <-throttle
            fmt.Println("+1")
            select {
            case req := <- sendSingleMessageRequests:
                fmt.Println("sending message", req)
                sendSingleMessage(req)
                fmt.Println("sent message", req)
            default:
                fmt.Println("no message sent")
            }
        }
    }()
}
