package solaredge

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
)

// HTTPError contains the HTTP error received from the SolarEdge API server
type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

var _ error = &HTTPError{}

// Error returns a string representation of an HTTPError
func (e *HTTPError) Error() string { return e.Status }

// Is return true if e2 is an HTTPError
func (e *HTTPError) Is(e2 error) bool {
	_, ok := e2.(*HTTPError)
	return ok
}

// ParseError wraps the error when failing to process the server response. Contains the original server response
// that generated the error, as well as the json error that triggered the error.
type ParseError struct {
	Err  error
	Body []byte
}

var _ error = &ParseError{}

// Error returns a string representation of a ParseError
func (e *ParseError) Error() string {
	return "json parse error: " + e.Err.Error()
}

// Unwrap returns the wrapped error
func (e *ParseError) Unwrap() error {
	return e.Err
}

// Is returns true if e2 is a ParseError
func (e *ParseError) Is(e2 error) bool {
	_, ok := e2.(*ParseError)
	return ok
}

// An APIError indicates the SolarEdge API server rejected the request (returning 403 - Forbidden).
//
// Apart from providing an invalid API key, this happens when the request arguments are not permitted,
// e.g. calling GetSiteEnergy() with a timeUnit of "DAY" and a time range of more than a year.
//
// WARNING: the SolarEdge API returns the error as an HTML document. APIError does its best to parse the document, but this may break at some point.
type APIError struct {
	Message string
}

var _ error = &APIError{}

// Error returns a string representation of an APIError
func (e *APIError) Error() string {
	return "api error: " + e.Message
}

// Is returns true if e2 is a ParseError
func (e *APIError) Is(e2 error) bool {
	_, ok := e2.(*APIError)
	return ok
}

func makeAPIError(body []byte) *APIError {
	if err := makeAPIErrorFromJSON(body); err != nil {
		return err
	}
	return makeAPIErrorFromHTML(body)
}

func makeAPIErrorFromJSON(body []byte) *APIError {
	var response struct {
		String string `json:"String"`
	}
	if err := json.Unmarshal(body, &response); err == nil {
		return &APIError{Message: response.String}
	}
	return nil
}

func makeAPIErrorFromHTML(body []byte) *APIError {
	var atMessage bool
	var message string

	tkn := html.NewTokenizer(bytes.NewBuffer(body))
loop:
	for {
		switch tkn.Next() {
		case html.ErrorToken:
			break loop
		case html.TextToken:
			tok := tkn.Token()
			if atMessage {
				message = tok.String()
				break loop
			}
			if tok.String() == "Message" {
				atMessage = true
			}
		}
	}

	if message == "" {
		message = "could not return error reason"
	}
	return &APIError{Message: message}
}
