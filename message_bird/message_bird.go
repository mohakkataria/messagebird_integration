package message_bird

import (
    "github.com/messagebird/go-rest-api"
    "github.com/mohakkataria/messagebird_integration/models"
    "github.com/mohakkataria/messagebird_integration/util"
    "time"
    "fmt"
)

var sendSingleMessageRequests chan models.SplitMessage
var messageBirdClient *messagebird.Client

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

func QueueMessage(message * models.Message) error {
    messageBodyLength := message.GetMessagebodyLength()

    splitMessage := message.GetSplitMessageWithOutBodyFromMessage()
    fmt.Println(splitMessage)
    // Here we  are making sure that the number of characters are as per specification.
    // Segmenting logic is also handled by dividing the message into chunks if required.
    if (!message.IsEncodingNormal()) {
        if (messageBodyLength < MAX_SINGLE_MSG_SEGMENT_SIZE_UNICODE) {
            splitMessage.Message = message.Message
            queueSplitMessage(splitMessage)
            fmt.Println(splitMessage)
        } else {
            if (messageBodyLength > MAX_MSG_SIZE_UNICODE) {
                message.Message = string([]rune(message.Message)[:MAX_MSG_SIZE_UNICODE])
            }

            messages := getMessageChunks(MAX_MSG_SIZE_UNICODE, MAX_MULTIPART_MSG_SEGMENT_SIZE_UNICODE, true, message.Message)
            segments := 0
            numberOfSegments := len(messages)
            referenceNumber := util.GetRandom2DigitHex()
            for _, body := range messages {
                segments++;
                splitMessage.Message = body
                /*
                    UDH format
                    Field 1 (1 octet): Length of User Data Header, in this case 05.
                    Field 2 (1 octet): Information Element Identifier, equal to 00
                    Field 3 (1 octet): Length of the header, excluding the first two fields; equal to 03
                    Field 4 (1 octet): 00-FF, CSMS reference number, must be same for all the SMS parts in the CSMS
                    Field 5 (1 octet): 00-FF, total number of parts.
                    Field 6 (1 octet): 00-FF, this part's number in the sequence. The value shall start at 1.
                */
                splitMessage.UDH = "050003" + referenceNumber + util.ConvertIntToHex(numberOfSegments) + util.ConvertIntToHex(segments);
                queueSplitMessage(splitMessage)
            }
        }
    } else {
        if (messageBodyLength < MAX_SINGLE_MSG_SEGMENT_SIZE_PLAIN) {
            splitMessage.Message = message.Message
            fmt.Println(splitMessage)
            queueSplitMessage(splitMessage)
        } else {
            // if message exceeds maximum size limit for plain messages, we only send the maximum allowed chars.
            if (messageBodyLength > MAX_MSG_SIZE_PLAIN) {
                message.Message = message.Message[:MAX_MULTIPART_MSG_SEGMENT_SIZE_PLAIN]
            }

            messages := getMessageChunks(MAX_MSG_SIZE_PLAIN, MAX_MULTIPART_MSG_SEGMENT_SIZE_PLAIN, false, message.Message)
            segments := 0
            numberOfSegments := len(messages)
            referenceNumber := util.GetRandom2DigitHex()
            for _, body := range messages {
                segments++;
                splitMessage.Message = body
                /*
                    UDH format
                    Field 1 (1 octet): Length of User Data Header, in this case 05.
                    Field 2 (1 octet): Information Element Identifier, equal to 00
                    Field 3 (1 octet): Length of the header, excluding the first two fields; equal to 03
                    Field 4 (1 octet): 00-FF, CSMS reference number, must be same for all the SMS parts in the CSMS
                    Field 5 (1 octet): 00-FF, total number of parts.
                    Field 6 (1 octet): 00-FF, this part's number in the sequence. The value shall start at 1.
                */
                splitMessage.UDH = "050003" + referenceNumber + util.ConvertIntToHex(numberOfSegments) + util.ConvertIntToHex(segments);
                queueSplitMessage(splitMessage)
            }
        }
    }
    return nil;
}

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

func queueSplitMessage(splitMessage models.SplitMessage) {
    sendSingleMessageRequests <- splitMessage
}

func sendSingleMessage(message models.SplitMessage) {
    // convert the SplitMessage to MessageParams and then retrieve the remaining properties and make a call to
    // message bird API.
    messageParams := &messagebird.MessageParams{TypeDetails:map[string]interface{}{"udh" : message.UDH}, DataCoding:message.DataCoding}
    mbMessage, err := messageBirdClient.NewMessage(message.Originator, message.Recipients, message.Message, messageParams)
    if err != nil {
        // messagebird.ErrResponse means custom JSON errors.
        if err == messagebird.ErrResponse {
            for _, mbError := range mbMessage.Errors {
                fmt.Printf("Error: %#v\n", mbError)
            }
        }
    }

}

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
