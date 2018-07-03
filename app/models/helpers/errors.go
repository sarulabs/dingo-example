package helpers

// ErrValidation is the error type that should be used
// to indicate that the error is caused by a validation problem.
type ErrValidation error

// ErrNotFound is the error type that should be used
// to indicate that the error is caused by the nonexistence of the requested resource.
type ErrNotFound error
