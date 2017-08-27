package models

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
        Recipients []int64
        Originator  string
        Message string
        Encoding Encoding
    }
)

