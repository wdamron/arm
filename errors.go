package arm

const (
	ErrInvalidInst     ErrorMessage = "invalid instruction id"
	ErrNoMatch         ErrorMessage = "no matching encoding"
	ErrInvalidEncoding ErrorMessage = "invalid instruction encoding"
)

type ErrorMessage string

func (err ErrorMessage) Error() string { return string(err) }
