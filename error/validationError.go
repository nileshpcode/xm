package error

import "fmt"

// NewValidationError creates a new standard validation error.
func NewValidationError(errorKey string, errors map[string]string) ValidationError {
	return ValidationError{ErrorKey: errorKey, Errors: errors}
}

// NewInvalidRequestPayloadError creates a new invalid request payload validation Error.
func NewInvalidRequestPayloadError(errorCode string) ValidationError {
	return ValidationError{ErrorKey: ErrorCodeInvalidRequestPayload, Errors: map[string]string{"payload": errorCode}}
}

// NewInvalidFieldsError creates a new invalid fields validation Error.
// 'failedFieldValidations' - map key should be the name of the field and value should be the error code.
func NewInvalidFieldsError(failedFieldValidations map[string]string) ValidationError {
	return ValidationError{ErrorKey: ErrorCodeInvalidFields, Errors: failedFieldValidations}
}

// ValidationError is an error indicating error in validations
type ValidationError struct {
	ErrorKey string            `json:"errorKey"`
	Errors   map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Error: [%s - %s]", e.ErrorKey, e.Errors)
}
