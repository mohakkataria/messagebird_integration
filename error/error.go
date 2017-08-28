// Package error contains description of Custom Error type
package error

const (
    BAD_JSON = "Bad Request JSON Input"
    BAD_RECIPIENT_INPUT = "Bad Recipient Input"
    MISSING_RECIPIENT_INPUT = "Missing Recipient Input"
    MISSING_ORIGINATOR_INPUT = "Missing Originator Input"
    ALPHANUMERIC_LENGTH_ORIGINATOR_ERROR = "Alphanumeric Originator should be less than 11 characters"
    BAD_ORIGINATOR_INPUT = "Bad Originator Input"
    NUMERIC_ORIGINATOR_ERROR = "Bad Originator Input as it cannot be a number less than 0"
    MISSING_MESSAGE_BODY = "Missing Message Body"
    BAD_MESSAGE_INPUT = "Bas Message Body Input"
)

type Error struct {
    Code int
    Message string
}
