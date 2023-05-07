package form3

// OperationError is used to provide a customized message that is easily consumable by the caller.
//
// It used while an operation is being executed and an error occurs.
type OperationError struct {
	Message string // Contains customized message, can contain the http status code if the http request was performed.
	Body    []byte // Contains the http body if the http request was performed.
}

// Error returns the message.
func (e OperationError) Error() string {
	return e.Message
}
