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

func (a NonRetryError) Error() string {
	return a.Msg
}
