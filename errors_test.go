package solaredge_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHTTPError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		other    error
		expectIs bool
		expectAs bool
	}{
		{
			name:     "direct",
			err:      &solaredge.HTTPError{StatusCode: http.StatusForbidden, Status: "Forbidden"},
			other:    &solaredge.HTTPError{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "wrapped",
			err:      fmt.Errorf("error: %w", &solaredge.HTTPError{StatusCode: http.StatusForbidden, Status: "Forbidden"}),
			other:    &solaredge.HTTPError{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "is not",
			err:      &solaredge.HTTPError{},
			other:    errors.New("error"),
			expectIs: false,
			expectAs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.err)
			assert.Equal(t, tt.expectIs, errors.Is(tt.err, tt.other))
			var newErr *solaredge.HTTPError
			assert.Equal(t, tt.expectAs, errors.As(tt.err, &newErr))
			if tt.expectAs {
				assert.ErrorIs(t, newErr, &solaredge.HTTPError{})
			}
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		other    error
		expectIs bool
		expectAs bool
	}{
		{
			name:     "direct",
			err:      &solaredge.ParseError{Err: &json.SyntaxError{Offset: 10}, Body: "hello"},
			other:    &solaredge.ParseError{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "wrapped",
			err:      fmt.Errorf("error: %w", &solaredge.ParseError{Err: &json.SyntaxError{Offset: 10}, Body: "hello"}),
			other:    &solaredge.ParseError{},
			expectIs: true,
			expectAs: true,
		},
		{
			name:     "is not",
			err:      &solaredge.ParseError{},
			other:    errors.New("error"),
			expectIs: false,
			expectAs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.err)
			assert.Equal(t, tt.expectIs, errors.Is(tt.err, tt.other))
			var newErr *solaredge.ParseError
			assert.Equal(t, tt.expectAs, errors.As(tt.err, &newErr))
			if tt.expectAs {
				assert.ErrorIs(t, newErr, &solaredge.ParseError{})
			}
		})
	}
}

func TestErrInvalidJSON_Unwrap(t *testing.T) {
	e := &solaredge.ParseError{
		Err: &json.SyntaxError{Offset: 10},
	}
	e2 := errors.Unwrap(e)
	var e3 *json.SyntaxError
	assert.ErrorAs(t, e2, &e3)
	assert.Equal(t, int64(10), e3.Offset)
}
