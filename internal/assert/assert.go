package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, value, expected T) {
	t.Helper()

	if value != expected {
		t.Errorf("expected equal; got: %v; want: %v", value, expected)
	}
}

func NotEqual[T comparable](t *testing.T, value, expected T) {
	t.Helper()

	if value == expected {
		t.Errorf("should not be: %v", value)
	}
}

func DeepEqual[T any](t *testing.T, value, expected T) {
	t.Helper()

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("expected deep equal; got: %v; want: %v", value, expected)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("expected error")
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("expected no error; got: %v", err)
	}
}

func ErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if !strings.Contains(err.Error(), expected) {
		t.Errorf("expected %s; got %v", expected, err.Error())
	}
}

func Nil(t *testing.T, value any) {
	t.Helper()

	if value != nil && !reflect.ValueOf(value).IsNil() {
		t.Errorf("expected nil; got %v", value)
	}
}

func NotNil(t *testing.T, value any) {
	t.Helper()

	if value == nil || reflect.ValueOf(value).IsNil() {
		t.Errorf("expected not nil")
	}
}

func StringContains(t *testing.T, value, expected string) {
	t.Helper()

	if !strings.Contains(value, expected) {
		t.Errorf("does not contain; got: %v; want to contain: %v", value, expected)
	}
}

func True(t *testing.T, actual bool) {
	t.Helper()

	if !actual {
		t.Errorf("expeted true")
	}
}
