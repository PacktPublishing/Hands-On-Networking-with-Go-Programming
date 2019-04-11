package gqlid

import "testing"

func TestInt64(t *testing.T) {
	tests := []struct {
		name             string
		id               *ID
		expectedService  string
		expectedResource string
		expectedInt64    int64
		expectedOK       bool
	}{
		{
			name:             "can roundtrip int64",
			id:               Int64("service_name", "resource_name", 123),
			expectedService:  "service_name",
			expectedResource: "resource_name",
			expectedInt64:    123,
			expectedOK:       true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual, ok := test.id.Int64()
			if test.expectedOK != ok {
				t.Fatalf("expected OK %v, got %v", test.expectedOK, ok)
			}
			if actual != test.expectedInt64 {
				t.Fatalf("expected %v, got %v", test.expectedInt64, actual)
			}
			encoded := test.id.Encoded()
			decoded, err := Parse(encoded)
			if err != nil {
				t.Fatalf("failed to encode and decode the ID: %v", err)
			}
			if decoded.ValueInt64 != actual {
				t.Errorf("failed to decode int64 value, expected %v, got %v", actual, decoded.ValueInt64)
			}
		})
	}
}
