package payloads

type (
    // Message struct payload is for API input description as defined in the problem statement
    // {"recipient":31612345678,"originator":"MessageBird","message":"This is a test message."}
    Message struct {
        Recipient string `json:recipient`
        Originator  string `json:originator`
        Message string `json:messageg`
    }

)

