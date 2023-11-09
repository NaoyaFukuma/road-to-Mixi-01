package testhelpers

import (
	"reflect"
	"strings"
	"testing"
)

// assertEqual is a test helper function to assert equality of two values.
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}

// AssertEqual checks if values are equal
func AssertDeepEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, got: %v", expected, actual)
	}
}

// assertNotEqual is a test helper function to assert inequality of two values.
func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected == actual {
		t.Errorf("Expected %v to not equal %v", expected, actual)
	}
}

// assertNil is a test helper function to assert that a value is nil.
func AssertNil(t *testing.T, actual interface{}) {
	t.Helper()
	if actual != nil {
		t.Errorf("Expected %v to be nil", actual)
	}
}

// assertNotNil is a test helper function to assert that a value is not nil.
func AssertNotNil(t *testing.T, actual interface{}) {
	t.Helper()
	if actual == nil {
		t.Errorf("Expected %v to not be nil", actual)
	}
}

// assertError is a test helper function to assert that an error is not nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}

// assertNoError is a test helper function to assert that an error is nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

// assertContains is a test helper function to assert that a string contains a substring.
func AssertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("Expected %q to contain %q", s, substr)
	}
}

// assertNotContains is a test helper function to assert that a string does not contain a substring.
func AssertNotContains(t *testing.T, s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("Expected %q to not contain %q", s, substr)
	}
}

// assertPanic is a test helper function to assert that a function panics.
func AssertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected function to panic, but it did not")
		}
	}()
	f()
}

// assertNotPanic is a test helper function to assert that a function does not panic.
func AssertNotPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected function to not panic, but it did")
		}
	}()
	f()
}

// assertPanicWith is a test helper function to assert that a function panics with a specific value.
func AssertPanicWith(t *testing.T, f func(), expected interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected function to panic, but it did not")
		} else if r != expected {
			t.Errorf("Expected function to panic with %v, but it panicked with %v", expected, r)
		}
	}()
	f()
}

// assertNotPanicWith is a test helper function to assert that a function does not panic with a specific value.
func AssertNotPanicWith(t *testing.T, f func(), expected interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Expected function to not panic, but it did")
		} else if r == expected {
			t.Errorf("Expected function to not panic with %v, but it did", expected)
		}
	}()
	f()
}
