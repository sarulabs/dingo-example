package helpers

// ErrValidation is the error type that should be used
// to indicate that the error is caused by a validation problem.
// The PublicMessage should be used as the error message for an external client.
type ErrValidation struct {
	Err           error
	PublicMessage string
}

// Error converts the ErrValidation to a string.
func (e ErrValidation) Error() string {
	return e.Err.Error()
}

// ErrNotFound is the error type that should be used
// to indicate that the error is caused by the nonexistence of the requested resource.
// The PublicMessage should be used as the error message for an external client.
type ErrNotFound struct {
	Err           error
	PublicMessage string
}

// Error converts the ErrNotFound to a string.
func (e ErrNotFound) Error() string {
	return e.Err.Error()
}

// ErrAlreadyExist is the error type that should be used
// when attempting to create a resource that already exists.
// The PublicMessage should be used as the error message for an external client.
type ErrAlreadyExist struct {
	Err           error
	PublicMessage string
}

// Error converts the ErrAlreadyExist to a string.
func (e ErrAlreadyExist) Error() string {
	return e.Err.Error()
}
