package webserver

import "net/http"

type RequestClient struct {
	http.Client
}

func newRequestClient() *RequestClient {
	return &RequestClient{
		*http.DefaultClient,
	}
}
