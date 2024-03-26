package internal

import (
	"os"
	"testing"
)

func TestReplacePrefix(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "Hello")
	param := Param{Name: "Test", Value: "${TEST_ENV_VAR}World"}
	value := param.ParseValue()
	os.Unsetenv("TEST_ENV_VAR")
	if value != "HelloWorld" {
		t.Fatalf("Expected %s to be HelloWorld", value)
	}
}

func TestReplacePostfix(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "World")
	param := Param{Name: "Test", Value: "Hello${TEST_ENV_VAR}"}
	value := param.ParseValue()
	os.Unsetenv("TEST_ENV_VAR")
	if value != "HelloWorld" {
		t.Fatalf("Expected %s to be HelloWorld", value)
	}
}

func TestReplaceMultiple(t *testing.T) {
	os.Setenv("TEST_VAR_1", "Hello")
	os.Setenv("TEST_VAR_2", "World")
	param := Param{Name: "Test", Value: "${TEST_VAR_1}${TEST_VAR_2}"}
	value := param.ParseValue()
	os.Unsetenv("TEST_VAR_1")
	os.Unsetenv("TEST_VAR_2")
	if value != "HelloWorld" {
		t.Fatalf("Expected %s to be HelloWorld", value)
	}
}
