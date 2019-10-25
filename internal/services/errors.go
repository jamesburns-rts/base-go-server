package services

type BadRequestError struct {
	Message string
	Cause   error
}

func InvalidInput(message string, errs ...error) *BadRequestError {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}

	return &BadRequestError{
		Message: message,
		Cause:   err,
	}
}

func (br *BadRequestError) BadRequestMessage() string {
	return br.Message
}
func (br *BadRequestError) Detail() error {
	return br.Cause
}
func (br *BadRequestError) Error() string {
	if br.Cause != nil {
		return errors.Wrap(br.Cause, br.Message).Error()
	}
	return br.Message
}
