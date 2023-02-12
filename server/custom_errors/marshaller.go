package custom_errors

import "fmt"

type MarshallerError struct {
	err error
}

func (m *MarshallerError) Error() string {
	return fmt.Sprintf("Error while marshalling, unmarshalling incoming request, err: %v", m.err)
}

func NewMarshallerError(err error) error {
	return &MarshallerError{err: err}
}
