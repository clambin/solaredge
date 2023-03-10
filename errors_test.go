package solaredge

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		other       error
		wantErrorIs func(assert.TestingT, error, error, ...interface{}) bool
	}{
		{
			name:        "direct",
			err:         &HTTPError{StatusCode: http.StatusForbidden, Status: "Forbidden"},
			other:       &HTTPError{},
			wantErrorIs: assert.ErrorIs,
		},
		{
			name:        "wrapped",
			err:         fmt.Errorf("error: %w", &HTTPError{StatusCode: http.StatusForbidden, Status: "Forbidden"}),
			other:       &HTTPError{},
			wantErrorIs: assert.ErrorIs,
		},
		{
			name:        "is not",
			err:         &HTTPError{},
			other:       errors.New("error"),
			wantErrorIs: assert.NotErrorIs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.err)
			tt.wantErrorIs(t, tt.err, tt.other)

			var httpError *HTTPError
			if errors.Is(tt.err, httpError) {
				assert.ErrorAs(t, tt.err, &httpError)
				assert.ErrorIs(t, httpError, &HTTPError{})
			}
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		other       error
		wantErrorIs func(assert.TestingT, error, error, ...interface{}) bool
	}{
		{
			name:        "direct",
			err:         &ParseError{Err: &json.SyntaxError{Offset: 10}, Body: []byte("hello")},
			other:       &ParseError{},
			wantErrorIs: assert.ErrorIs,
		},
		{
			name:        "wrapped",
			err:         fmt.Errorf("error: %w", &ParseError{Err: &json.SyntaxError{Offset: 10}, Body: []byte("hello")}),
			other:       &ParseError{},
			wantErrorIs: assert.ErrorIs,
		},
		{
			name:        "is not",
			err:         &ParseError{},
			other:       errors.New("error"),
			wantErrorIs: assert.NotErrorIs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.err)
			tt.wantErrorIs(t, tt.err, tt.other)

			var newErr *ParseError
			if errors.Is(tt.err, newErr) {
				assert.ErrorAs(t, tt.err, &newErr)
				assert.ErrorIs(t, newErr, &ParseError{})
			}
		})
	}
}

func TestErrInvalidJSON_Unwrap(t *testing.T) {
	e := &ParseError{
		Err: &json.SyntaxError{Offset: 10},
	}
	e2 := errors.Unwrap(e)
	var e3 *json.SyntaxError
	assert.ErrorAs(t, e2, &e3)
	assert.Equal(t, int64(10), e3.Offset)
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid - json",
			input:    `{ "String": "invalid token" }`,
			expected: "invalid token",
		},
		{
			name:     "valid - html",
			input:    `<!doctype html><html lang="en"><head><title>HTTP Status 403 – Forbidden</title><style type="text/css">h1 {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;font-size:22px;} h2 {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;font-size:16px;} h3 {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;font-size:14px;} body {font-family:Tahoma,Arial,sans-serif;color:black;background-color:white;} b {font-family:Tahoma,Arial,sans-serif;color:white;background-color:#525D76;} p {font-family:Tahoma,Arial,sans-serif;background:white;color:black;font-size:12px;} a {color:black;} a.name {color:black;} .line {height:1px;background-color:#525D76;border:none;}</style></head><body><h1>HTTP Status 403 – Forbidden</h1><hr class="line" /><p><b>Type</b> Status Report</p><p><b>Message</b> The requested time frame exceed the allowed limit for requested resource</p><p><b>Description</b> The server understood the request but refuses to authorize it.</p><hr class="line" /><h3>Apache Tomcat/8.5.46</h3></body></html>`,
			expected: " The requested time frame exceed the allowed limit for requested resource",
		},
		{
			name:     "empty",
			input:    ``,
			expected: "could not return error reason",
		},
		{
			name:     "invalid",
			input:    `not a valid HTML document`,
			expected: "could not return error reason",
		},
		{
			name:     "missing message",
			input:    `<html><head><title>HTTP Status 403 – Forbidden</title></head><body><h1>HTTP Status 403 – Forbidden</h1><hr/><p><b>Type</b> Status Report</p><p><b>Description</b> The server understood the request but refuses to authorize it.</p><hr class="line" /><h3>Apache Tomcat/8.5.46</h3></body></html>`,
			expected: "could not return error reason",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeAPIError([]byte(tt.input))
			require.Error(t, err)
			assert.ErrorIs(t, err, &APIError{})
			assert.Equal(t, "api error: "+tt.expected, err.Error())
		})
	}
}
