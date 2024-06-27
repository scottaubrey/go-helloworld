package api

type RequestError struct {
	Url      string
	HTTPCode int
	Body     string
	Err      string
}

func (requestError RequestError) Error() string {
	return requestError.Err
}
