package main

import (
	"testing"
)

func TestShardSection(t *testing.T) {
	tests := []struct {
		name       string
		shardCount int
		input      byte
		expected   int
	}{
		{
			name:       "zero value is 0",
			shardCount: 32,
			input:      0,
			expected:   0,
		},
		{
			name:       "1 shards to 0",
			shardCount: 32,
			input:      1,
			expected:   0,
		},
		{
			name:       "8 shards to 1",
			shardCount: 32,
			input:      8,
			expected:   1,
		},
		{
			name:       "247 shards to 1",
			shardCount: 32,
			input:      255,
			expected:   31,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s, err := NewSharder(tc.shardCount)
			if err != nil {
				t.Fatalf("failed to create sharder: %v", err)
			}
			actual := s.shardIndex(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}
