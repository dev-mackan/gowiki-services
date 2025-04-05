package webserver

type WebErrorResponse interface {
	Error() string         // Response to client
	InternalError() string // The actual error
	Status() int           // http status code
}

type FetchError struct {
	InternalError string
}

func (e *FetchError) Error() string {
	return "Error fetching a resource."
}
func (e *FetchError) Status() int {
	return 500
}

type ReadResourceError struct{}

func (e *ReadResourceError) Error() string {
	return "Error reading a response"
}

func (e *ReadResourceError) Status() int {
	return 500
}

type UnmarshalResourceError struct{}

func (e *UnmarshalResourceError) Error() string {
	return "Error parsing a resource"
}

func (e *UnmarshalResourceError) Status() int {
	return 500
}
