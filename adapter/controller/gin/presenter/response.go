package presenter

func NewErrorResponse(code int, message string) (int, *ErrorResponse) {
	return code, &ErrorResponse{
		Error: Error{
			Message: message,
			Code:    code,
		},
	}
}
