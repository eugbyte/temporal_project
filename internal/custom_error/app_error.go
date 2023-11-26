package customerrors

// To indicate to temporal to not retry on error
type ApplicationError struct {
	Msg string
}

func NewAppError(msg string) ApplicationError {
	return ApplicationError{
		Msg: msg,
	}
}

func (a ApplicationError) Error() string {
	return a.Msg
}
