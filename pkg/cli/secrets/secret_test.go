package secrets

import (
	"testing"
)

func TestSecret_DestroyValue(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
	}{
		{
			name:  "nil value",
			value: nil,
		},
		{
			name:  "empty value",
			value: []byte{},
		},
		{
			name:  "non-empty value",
			value: []byte("secret-data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a secret with the test value
			s := NewSecret(tt.value)
			
			// Call DestroyValue
			s.DestroyValue()
			
			// Check that Value is nil
			if s.Value != nil {
				t.Errorf("DestroyValue() did not set Value to nil, got %v", s.Value)
			}
		})
	}
}

func TestSecret_DestroyValue_ZeroesMemory(t *testing.T) {
	// Create a secret with a non-empty value
	originalValue := []byte("sensitive-data")
	s := NewSecret(make([]byte, len(originalValue)))
	copy(s.Value, originalValue)
	
	// Get reference to the value slice for later inspection
	valueRef := s.Value
	
	// Call DestroyValue
	s.DestroyValue()
	
	// Verify the original memory was zeroed before being set to nil
	for i, b := range valueRef {
		if b != 0 {
			t.Errorf("DestroyValue() did not zero memory at index %d, got %d", i, b)
		}
	}
}