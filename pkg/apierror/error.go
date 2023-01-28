package apierror

type ApiError struct {
}

func (err ApiError) Error() string {
	return "hardcoded"
}
