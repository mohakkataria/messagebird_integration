// Package models contains the models and associated methods if any for Messages
package models

import "unicode/utf8"

const (
	// NORMAL Encoding representation in Enum
	NORMAL Encoding = iota // 0
	// UNICODE Encoding representation in Enum
	UNICODE // 1
)

type (
	// Encoding as Enum
	Encoding int

	// Message struct model is for description of input as defined in the problem statement
	// {"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}
	Message struct {
		Recipients  []string
		Originator  string
		MessageBody string
		Encoding    Encoding
	}

	// SplitMessage struct model is for Concatenated SMS used in queueing the requests before sending them
	SplitMessage struct {
		Recipients       []string
		Originator       string
		MessageBodyChunk string
		DataCoding       string
		UDH              string
	}
)


// IsEncodingNormal is used to check where encoding of the Message is Normal or not
func (m Message) IsEncodingNormal() bool {
	return m.Encoding == NORMAL
}

// GetMessagebodyLength returns the body length regardless of encoding.
func (m Message) GetMessagebodyLength() int {
	if !m.IsEncodingNormal() {
		return utf8.RuneCountInString(m.MessageBody)
	}
	return len(m.MessageBody)
}

// GetSplitMessageWithOutBodyFromMessage returns a SplitMessage from Message with metadata of udh.
func (m Message) GetSplitMessageWithOutBodyFromMessage() SplitMessage {
	splitMessage := SplitMessage{}
	splitMessage.Recipients = m.Recipients
	splitMessage.Originator = m.Originator
	if m.IsEncodingNormal() {
		splitMessage.DataCoding = "plain"
	} else {
		splitMessage.DataCoding = "unicode"
	}
	return splitMessage
}
