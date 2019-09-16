package verrors

// InternalError Database error and etc.
type InternalError struct {
	Message string `json:"message"`
}

func (e InternalError) Error() string {
	return e.Message
}