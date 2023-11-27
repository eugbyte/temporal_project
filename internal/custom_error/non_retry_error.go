package customerrors

// To indicate to temporal to not retry on error
type NonRetryError struct {
	Msg string
}

func NewNonRetryError(msg string) NonRetryError {
	return NonRetryError{
		Msg: msg,
	}
}

func (e NonRetryError) Error() string {
	return e.Msg
}
