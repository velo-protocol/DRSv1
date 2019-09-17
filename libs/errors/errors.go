package verrors

type ParameterError struct {
	Message string `json:"message"`
}

func (e ParameterError) Error() string {
	return e.Message
}

type InternalError struct {
	Message string `json:"message"`
}

func (e InternalError) Error() string {
	return e.Message
}