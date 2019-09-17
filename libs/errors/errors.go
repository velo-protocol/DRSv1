package verrors

// ParameterError Invalid parameter
type ParameterError struct {
	Message string `json:"message"`
}

func (e ParameterError) Error() string {
	return e.Message
}

// InternalError Database error and etc.
type InternalError struct {
	Message string `json:"message"`
}

func (e InternalError) Error() string {
	return e.Message
}