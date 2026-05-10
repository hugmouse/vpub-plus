package model

import "testing"

func TestNormalizeRestrictedVisibility(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RestrictedVisibility
	}{
		{
			name:     "visible stays visible",
			input:    "visible",
			expected: RestrictedVisibilityVisible,
		},
		{
			name:     "hidden stays hidden",
			input:    "hidden",
			expected: RestrictedVisibilityHidden,
		},
		{
			name:     "empty defaults to hidden",
			input:    "",
			expected: RestrictedVisibilityHidden,
		},
		{
			name:     "invalid defaults to hidden",
			input:    "invalid",
			expected: RestrictedVisibilityHidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeRestrictedVisibility(tt.input); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
