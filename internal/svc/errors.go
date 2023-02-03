package svc

type InvalidInputError struct {
	Message string
	Cause   error
}

func (i InvalidInputError) Error() string {
	return i.Message
}

type BadRequestError struct {
	Message string
	Cause   error
}

func (i BadRequestError) Error() string {
	return i.Message
}
