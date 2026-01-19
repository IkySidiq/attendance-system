package exceptions

type HTTPError interface {
	error
	StatusCode() int
}	

type ClientError struct {
	code    int
	message string
}

func NewClientError(code int, message string) *ClientError {
	return &ClientError{
		code:    code,
		message: message,
	}
}

func (e *ClientError) Error() string {
	return e.message
}

func (e *ClientError) StatusCode() int {
	return e.code
}