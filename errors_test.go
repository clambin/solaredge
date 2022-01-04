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
	err := &solaredge.HTTPError{
		StatusCode: http.StatusForbidden,
		Status:     "403 Forbidden",
	}

	assert.Equal(t, "403 Forbidden", err.Error())
	assert.True(t, errors.Is(err, &solaredge.HTTPError{}))

	err2 := fmt.Errorf("http: %w", err)
	assert.Equal(t, "http: 403 Forbidden", err2.Error())
	assert.True(t, errors.Is(err, &solaredge.HTTPError{}))

	err3 := &solaredge.HTTPError{}
	assert.True(t, errors.As(err, &err3))
	assert.Equal(t, "403 Forbidden", err3.Error())
}

func TestParseError(t *testing.T) {
	err := &solaredge.ParseError{
		Body: "foo",
		Err:  &json.SyntaxError{Offset: 10},
	}

	assert.Equal(t, "json parse error: ", err.Error())
	assert.True(t, errors.Is(err, &solaredge.ParseError{}))

	err2 := fmt.Errorf("error: %w", err)
	assert.Equal(t, "error: json parse error: ", err2.Error())
	assert.True(t, errors.Is(err, &solaredge.ParseError{}))

	err3 := &solaredge.ParseError{}
	assert.True(t, errors.As(err2, &err3))
	assert.Equal(t, "json parse error: ", err3.Error())
	assert.Equal(t, "foo", err3.Body)

	err4 := &json.SyntaxError{}
	assert.True(t, errors.As(err2, &err4))
	assert.Equal(t, int64(10), err4.Offset)
}
