package containerator

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, actual, expected interface{}, message string) {
	if expected == nil {
		check := actual == nil
		if !check {
			val := reflect.ValueOf(actual)
			check = val.Kind() == reflect.Ptr && val.IsNil()
		}
		if !check {
			t.Fatalf("%s - got: %v / want: %v", message, actual, expected)
		}
		return
	}
	if actual != expected {
		t.Fatalf("%s - got: %v / want: %v", message, actual, expected)
	}
}
