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
	sc := http.StatusForbidden
	s := http.StatusText(sc)

	err := &solaredge.HTTPError{
		StatusCode: sc,
		Status:     s,
	}

	assert.Equal(t, s, err.Error())
	assert.True(t, errors.Is(err, &solaredge.HTTPError{}))

	err2 := fmt.Errorf("http: %w", err)
	assert.True(t, errors.Is(err2, &solaredge.HTTPError{}))
	assert.Equal(t, "http: "+s, err2.Error())

	err3 := &solaredge.HTTPError{}
	assert.True(t, errors.As(err2, &err3))
	assert.Equal(t, s, err3.Error())
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
