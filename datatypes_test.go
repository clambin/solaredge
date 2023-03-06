package solaredge

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     Date
		wantErr  assert.ErrorAssertionFunc
		stringed string
	}{
		{
			name:     "valid date",
			input:    `"2023-03-04"`,
			want:     Date(time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)),
			wantErr:  assert.NoError,
			stringed: "2023-03-04 00:00:00 +0000 UTC",
		},
		{
			name:    "invalid date",
			input:   `"2023-03-04 14:00:00"`,
			wantErr: assert.Error,
		},
		{
			name:    "blank",
			input:   `""`,
			wantErr: assert.Error,
		},
		{
			name:    "empty",
			input:   ``,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalled Date
			err := json.Unmarshal([]byte(tt.input), &unmarshalled)

			assert.Equal(t, tt.want, unmarshalled)
			tt.wantErr(t, err)

			if time.Time(unmarshalled).IsZero() {
				return
			}

			assert.Equal(t, tt.stringed, unmarshalled.String())
			marshalled, err := json.Marshal(unmarshalled)
			require.NoError(t, err)
			assert.Equal(t, tt.input, string(marshalled))
		})
	}
}

func TestTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     Time
		wantErr  assert.ErrorAssertionFunc
		stringed string
	}{
		{
			name:     "valid date",
			input:    `"2023-03-04 14:30:15"`,
			want:     Time(time.Date(2023, time.March, 4, 14, 30, 15, 0, time.UTC)),
			wantErr:  assert.NoError,
			stringed: "2023-03-04 14:30:15 +0000 UTC",
		},
		{
			name:    "invalid date",
			input:   `"2023-03-04"`,
			wantErr: assert.Error,
		},
		{
			name:    "blank",
			input:   `""`,
			wantErr: assert.Error,
		},
		{
			name:    "empty",
			input:   ``,
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalled Time
			err := json.Unmarshal([]byte(tt.input), &unmarshalled)

			assert.Equal(t, tt.want, unmarshalled)
			tt.wantErr(t, err)

			if time.Time(unmarshalled).IsZero() {
				return
			}

			assert.Equal(t, tt.stringed, unmarshalled.String())

			marshalled, err := json.Marshal(unmarshalled)
			assert.Equal(t, tt.input, string(marshalled))
			assert.NoError(t, err)
		})
	}
}
