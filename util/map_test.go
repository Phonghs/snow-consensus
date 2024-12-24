package util

import (
	"reflect"
	"testing"
)

func TestGetKeys(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]int
		expected []string
	}{
		{
			name:     "Empty map",
			inputMap: map[string]int{},
			expected: []string{},
		},
		{
			name:     "Nil map",
			inputMap: nil,
			expected: []string{},
		},
		{
			name: "Map with keys",
			inputMap: map[string]int{
				"key1": 1,
				"key2": 2,
				"key3": 3,
			},
			expected: []string{"key1", "key2", "key3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := GetKeys(tt.inputMap)

			if !reflect.DeepEqual(keys, tt.expected) {
				t.Errorf("GetKeys(%v) = %v, expected %v", tt.inputMap, keys, tt.expected)
			}
		})
	}
}

func TestGetMapExcludingKey(t *testing.T) {
	tests := []struct {
		name        string
		originalMap map[string]int
		excludeKey  string
		expected    map[string]int
	}{
		{
			name:        "Empty map",
			originalMap: map[string]int{},
			excludeKey:  "key1",
			expected:    map[string]int{},
		},
		{
			name: "Map with one key",
			originalMap: map[string]int{
				"key1": 1,
			},
			excludeKey: "key1",
			expected:   map[string]int{},
		},
		{
			name: "Map excluding a key",
			originalMap: map[string]int{
				"key1": 1,
				"key2": 2,
				"key3": 3,
			},
			excludeKey: "key2",
			expected: map[string]int{
				"key1": 1,
				"key3": 3,
			},
		},
		{
			name: "Exclude non-existing key",
			originalMap: map[string]int{
				"key1": 1,
				"key2": 2,
			},
			excludeKey: "key3",
			expected: map[string]int{
				"key1": 1,
				"key2": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMapExcludingKey(tt.originalMap, tt.excludeKey)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetMapExcludingKey(%v, %v) = %v, expected %v", tt.originalMap, tt.excludeKey, result, tt.expected)
			}
		})
	}
}
