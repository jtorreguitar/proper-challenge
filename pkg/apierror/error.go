package apierror

type ApiError struct {
	Code       string
	InnerCause error
	Values     map[string]any
}

func (err ApiError) Error() string {
	if err.InnerCause != nil {
		return err.Code + ": " + err.InnerCause.Error()
	}

	return err.Code
}
