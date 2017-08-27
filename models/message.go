package models

type (
    // Message struct model is for description as defined in the problem statement
    // {"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}
    Message struct {
        recipient int64 `json:recipient`
        originator  string `json:originator`
        message string `json:messageg`
    }
)
