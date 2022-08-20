package arm

const (
	ErrInvalidInst     ErrorMessage = "invalid instruction id"
	ErrNoMatch         ErrorMessage = "no matching encoding"
	ErrInvalidEncoding ErrorMessage = "invalid instruction encoding"
)

// ErrorMessage is an error message type, returned when instruction matching or encoding fails.
type ErrorMessage string

func (err ErrorMessage) Error() string { return string(err) }
