// Package error contains description of Custom Error type
package error

// Constants for different error strings used
const (
	BadJSON = "Bad Request JSON Input"
	BadRecipientInput = "Bad Recipient Input"
	MissingRecipientInput = "Missing Recipient Input"
	MissingOriginatorInput = "Missing Originator Input"
	AlphanumericLengthOriginatorError = "Alphanumeric Originator should be less than 11 characters"
	BadOriginatorInput = "Bad Originator Input"
	NumericOriginatorInput = "Bad Originator Input as it cannot be a number less than 0"
	MissingMessageBody = "Missing Message Body"
	BadMessageInput = "Bas Message Body Input"
)

// Error type is a description of the custom error type with code and message fields
type Error struct {
	Code    int
	Message string
}
