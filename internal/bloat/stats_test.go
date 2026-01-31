package bloat

import (
	"testing"
	"time"
)

func TestFormatInterval(t *testing.T) {
	tests := []struct {
		name     string
		value    time.Duration
		expected string
	}{
		{name: "zero", value: 0, expected: "0 seconds"},
		{name: "minutes", value: 20 * time.Minute, expected: "20 minutes"},
		{name: "hours", value: 3 * time.Hour, expected: "3 hours"},
		{name: "days", value: 48 * time.Hour, expected: "2 days"},
		{name: "sub-minute", value: 30 * time.Second, expected: "1 minute"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatInterval(tt.value)
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
