package containerator

import "testing"

func assertEqual(t *testing.T, actual, expected interface{}, message string) {
	if actual != expected {
		t.Fatalf("%s - got: %v / want: %v", message, actual, expected)
	}
}
