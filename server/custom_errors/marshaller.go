package custom_errors

import "fmt"

/**
Everything here are custom error objects we use to dynamically parse and generate what statusCode to return to
front-end.

Most errors here are self-explanatory.
*/

type MarshallerError struct {
	err error
}

func (m *MarshallerError) Error() string {
	return fmt.Sprintf("Error while marshalling, unmarshalling incoming request, err: %v", m.err)
}

func NewMarshallerError(err error) error {
	return &MarshallerError{err: err}
}
