package main

import (
	"fmt"
	"net/http"
)

type MyJwtTransport struct {
	transport http.RoundTripper
	token     string
}

func (t MyJwtTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.token != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	}
	return t.transport.RoundTrip(request)
}
