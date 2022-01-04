package solaredge

// HTTPError contains the HTTP error received from the SolarEdge API server
type HTTPError struct {
	StatusCode int
	Status     string
}

var _ error = &HTTPError{}

// Error call complies with the error interface
func (e *HTTPError) Error() string { return e.Status }

// Is checks if a received error is an HTTPError
func (e *HTTPError) Is(e2 error) bool {
	_, ok := e2.(*HTTPError)
	return ok
}

// ParseError wraps the error when failing to process the server response. Contains the original server response
// that generated the error, as well as the json error that triggered the error.
type ParseError struct {
	Body string
	Err  error
}

var _ error = &ParseError{}

// Error call complies with the error interface
func (e *ParseError) Error() string {
	return "json parse error: " + e.Err.Error()
}

// Unwrap returns the wrapped error
func (e *ParseError) Unwrap() error {
	return e.Err
}

// Is checks if a received error is a ParseError
func (e *ParseError) Is(e2 error) bool {
	_, ok := e2.(*ParseError)
	return ok
}
