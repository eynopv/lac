package param

import (
	"os"
	"testing"
)

func TestReplacePrefixWithEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "Hello")
	param := Param("${TEST_ENV_VAR}World")
	expected := "HelloWorld"
	result := param.Resolve(nil)
	os.Unsetenv("TEST_ENV_VAR")
	if result != expected {
		t.Fatalf("Expected %s to be HelloWorld", result)
	}
}

func TestReplacePostfixWithEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "World")
	param := Param("Hello${TEST_ENV_VAR}")
	expected := "HelloWorld"
	result := param.Resolve(nil)
	os.Unsetenv("TEST_ENV_VAR")
	if expected != "HelloWorld" {
		t.Fatalf("Expected %s to be HelloWorld", result)
	}
}

func TestReplaceMultipleWithEnv(t *testing.T) {
	os.Setenv("TEST_VAR_1", "Hello")
	os.Setenv("TEST_VAR_2", "World")
	param := Param("${TEST_VAR_1}${TEST_VAR_2}")
	expected := "HelloWorld"
	result := param.Resolve(nil)
	os.Unsetenv("TEST_VAR_1")
	os.Unsetenv("TEST_VAR_2")
	if result != expected {
		t.Fatalf("Expected %s to be HelloWorld", result)
	}
}

func TestReplaceMultipleWithReplacements(t *testing.T) {
	param := Param("${replacement_1}, ${replacement_2}!")
	replacements := map[string]string{
		"replacement_1": "Hello",
		"replacement_2": "World",
	}
	expected := "Hello, World!"
	result := param.Resolve(replacements)
	if result != expected {
		t.Fatalf("Expected %s to be HelloWorld", result)
	}
}
