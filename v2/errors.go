package v2

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

type Error struct {
	vals map[string]any
}

func (e *Error) Error() string {
	return fmt.Sprintf("api error: %v", e.vals)
}

func newResponseError(r *http.Response) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		// parse json error
		var e map[string]any
		if err := json.NewDecoder(r.Body).Decode(&e); err == nil {
			return &Error{vals: e}
		}
	case "text/html":
		if values := readHTMLError(r.Body); len(values) > 0 {
			return &Error{vals: values}
		}
	}
	// if response is neither (valid) json nor html, report the HTTP Status Code (& Status, if available).
	if r.Status != "" {
		return fmt.Errorf("api error: %s", r.Status)
	}
	return fmt.Errorf("api error: %d - %s", r.StatusCode, http.StatusText(r.StatusCode))
}

func readHTMLError(r io.Reader) map[string]any {
	values := make(map[string]any)
	tokenizer := html.NewTokenizer(r)

	var currentKey string

	for {
		// Stop parsing if both keys are found
		if len(values) == 2 {
			return values
		}

		switch tokenizer.Next() {
		case html.ErrorToken:
			// Return EOF as nil error for clean termination
			if tokenizer.Err() == io.EOF {
				return values
			}
			return nil

		case html.TextToken:
			text := /*html.UnescapeString(*/ string(tokenizer.Text()) //)

			// Check if we're currently expecting a value for a key
			if currentKey != "" {
				values[currentKey] = text
				currentKey = "" // Reset the key after storing the value
			} else if text == "Message" {
				currentKey = "message"
			} else if text == "Description" {
				currentKey = "description"
			}

		default:
			// Skip non-text tokens
			continue
		}
	}
}
