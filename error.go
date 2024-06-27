package main

type RequestError struct {
	HTTPCode int
	Body     string
	Err      string
}

func (requestError RequestError) Error() string {
	return requestError.Err
}
