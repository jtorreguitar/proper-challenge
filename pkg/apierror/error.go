package apierror

type ApiError struct {
	Code       string
	InnerCause error
}

func (err ApiError) Error() string {
	return err.Code + ": " + err.InnerCause.Error()
}
