package utils
type HttpError struct {
	Err string `json:"err"`
}

func (h HttpError) Error() string {
	return h.Err
}

