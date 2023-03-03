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
		pass     bool
		expected Date
		stringed string
	}{
		{
			name:     "valid date",
			input:    `"2023-03-04"`,
			pass:     true,
			expected: Date(time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)),
			stringed: "2023-03-04 00:00:00 +0000 UTC",
		},
		{
			name:  "invalid date",
			input: `"2023-03-04 14:00:00"`,
			pass:  false,
		},
		{
			name:  "blank",
			input: `""`,
			pass:  false,
		},
		{
			name:  "empty",
			input: ``,
			pass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalled Date
			err := json.Unmarshal([]byte(tt.input), &unmarshalled)
			if !tt.pass {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.expected, unmarshalled)

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
		pass     bool
		expected Time
		stringed string
	}{
		{
			name:     "valid date",
			input:    `"2023-03-04 14:30:15"`,
			pass:     true,
			expected: Time(time.Date(2023, time.March, 4, 14, 30, 15, 0, time.UTC)),
			stringed: "2023-03-04 14:30:15 +0000 UTC",
		},
		{
			name:  "invalid date",
			input: `"2023-03-04"`,
			pass:  false,
		},
		{
			name:  "blank",
			input: `""`,
			pass:  false,
		},
		{
			name:  "empty",
			input: ``,
			pass:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unmarshalled Time
			err := json.Unmarshal([]byte(tt.input), &unmarshalled)
			if !tt.pass {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.expected, unmarshalled)

			assert.Equal(t, tt.stringed, unmarshalled.String())

			marshalled, err := json.Marshal(unmarshalled)
			require.NoError(t, err)
			assert.Equal(t, tt.input, string(marshalled))
		})
	}
}
