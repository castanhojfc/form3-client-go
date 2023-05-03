package form3

import (
	"fmt"
)

// OperationError should be used to provide a customized message when an undesired response is obtained from the API
type OperationError struct {
	Message  string // Custom message to help the user understand exactly what operation failed
	Response string // Full HTTP response details provided by the API which originated the anomalous behaviour
}

func (e OperationError) Error() string {
	return fmt.Sprintf("%s, Response: %s", e.Message, e.Response)
}
