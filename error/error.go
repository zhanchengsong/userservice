package error

type Error interface {
	getCode() int
	getMessage() string
}

type DBError struct {
	code int
	message string
}

func (dbe DBError) getCode() int{
	return dbe.code
}

func (dbe DBError) getMessage() string{
	return dbe.message
}