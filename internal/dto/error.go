package dto

// Error represents a generic error message
type Error struct {
	Message string
}

func (e *Error) String() string {
	return e.Message
}
