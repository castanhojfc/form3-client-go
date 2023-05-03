package form3

import (
	"fmt"
)

type OperationError struct {
	Message  string
	Response string
}

func (e OperationError) Error() string {
	return fmt.Sprintf("%s, Response: %s", e.Message, e.Response)
}
