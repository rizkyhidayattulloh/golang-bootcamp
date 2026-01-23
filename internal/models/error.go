package models

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func NewError(status int, message string) *Error {
	return &Error{
		Status:  status,
		Message: message,
	}
}
