package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseStringUnixMilliSecTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		expected  time.Time
	}{
		{
			name:      "Unix timestamp in seconds (integer)",
			timestamp: "1609459200", // Corresponds to 2021-01-01 00:00:00 UTC
			expected:  time.Unix(1609459200, 0),
		},
		{
			name:      "Unix timestamp in milliseconds (13 digits)",
			timestamp: "1609459200123", // Corresponds to 2021-01-01 00:00:00.123 UTC
			expected:  time.Unix(1609459200, 123*1e6), // 123 milliseconds
		},
		{
			name:      "Unix timestamp with sub-second precision (float)",
			timestamp: "1609459200.123", // Corresponds to 2021-01-01 00:00:00.123 UTC
			expected:  time.Unix(1609459200, 123*1e6), // 123 milliseconds
		},
		{
			name:      "Unix timestamp with more sub-second precision",
			timestamp: "1609459200.123456789", // Should round to nanoseconds
			expected:  time.Unix(1609459200, 123456789), // 123456789 nanoseconds
		},
		{
			name:      "Invalid timestamp falls back to time.Now()",
			timestamp: "invalid", // Invalid input
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseStringUnixMilliSecTimestamp(tt.timestamp)

			if tt.timestamp == "invalid" {
				// We can't predict time.Now(), so just ensure that a valid time is returned
				assert.NotZero(t, actual, "Expected time.Now(), but got zero time")
			} else {
				// Use assert.InDelta to allow a small delta in the comparison to handle floating point precision issues
				assert.InDelta(t, tt.expected.UnixNano(), actual.UnixNano(), 1000,
					"Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
