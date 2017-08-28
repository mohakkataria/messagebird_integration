// Package models contains the models and associated methods if any for Messages
package models

import "unicode/utf8"

const (
    // Normal Encoding representation in Enum
    NORMAL Encoding = iota // 0
    // Unicode Encoding representation in Enum
    UNICODE                // 1
)

type (
    // Encoding as Enum
    Encoding int

    // Message struct model is for description as defined in the problem statement
    // {"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}
    Message struct {
        Recipients  []string
        Originator  string
        MessageBody string
        Encoding    Encoding
    }

    SplitMessage struct {
        Recipients       []string
        Originator       string
        MessageBodyChunk string
        DataCoding       string
        UDH              string
    }
)

func (this Message) IsEncodingNormal() bool {
    return this.Encoding == NORMAL
}

func (this Message) GetMessagebodyLength() int {
    if !this.IsEncodingNormal() {
        return utf8.RuneCountInString(this.MessageBody)
    }
    return len(this.MessageBody)
}

func (this Message) GetSplitMessageWithOutBodyFromMessage() SplitMessage {
    splitMessage := SplitMessage{}
    splitMessage.Recipients = this.Recipients
    splitMessage.Originator = this.Originator
    if this.IsEncodingNormal() {
        splitMessage.DataCoding = "plain"
    } else {
        splitMessage.DataCoding = "unicode"
    }
    return splitMessage
}

