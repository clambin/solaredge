package solaredge

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name string
		resp http.Response
		want string
	}{
		{
			name: "json error",
			resp: http.Response{
				Header: http.Header{"Content-Type": []string{"application/json"}},
				Body:   io.NopCloser(strings.NewReader(`{"message": "error"}`)),
			},
			want: `api error: map[message:error]`,
		},
		{
			name: "invalid json: use http error",
			resp: http.Response{
				Status: "bad request",
				Header: http.Header{"Content-Type": []string{"application/json"}},
				Body:   io.NopCloser(strings.NewReader(``)),
			},
			want: `api error: bad request`,
		},
		{
			name: "html error",
			resp: http.Response{
				Header: http.Header{"Content-Type": []string{"text/html"}},
				Body: io.NopCloser(strings.NewReader(`
<html>
<body>
	<h1>HTTP Status 400 – Bad Request</h1>
	<hr class="line" />
	<p><b>Type</b> Status Report</p>
	<p><b>Message</b> Required parameter &#39;startDate&#39; is not present</p>
	<p><b>Description</b> required parameter missing.</p>
	<hr class="line" />
	<h3>Apache Tomcat/8.5.46</h3>
</body>
</html>`)),
			},
			want: `api error: map[description: required parameter missing. message: Required parameter 'startDate' is not present]`,
		},
		{
			name: "invalid html: use http error",
			resp: http.Response{
				Status: "bad request",
				Header: http.Header{"Content-Type": []string{"text/html"}},
				Body: io.NopCloser(strings.NewReader(`
<html>
`)),
			},
			want: `api error: bad request`,
		},
		{
			name: "http error (no status)",
			resp: http.Response{StatusCode: http.StatusBadRequest},
			want: `api error: 400 - Bad Request`,
		},
		{
			name: "http error",
			resp: http.Response{StatusCode: http.StatusBadRequest, Status: http.StatusText(http.StatusBadRequest)},
			want: `api error: Bad Request`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newResponseError(&tt.resp)
			if e.Error() != tt.want {
				t.Errorf("got %q, want %q", e.Error(), tt.want)
			}
		})
	}
}
