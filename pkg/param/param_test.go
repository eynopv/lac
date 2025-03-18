package param

import (
	"os"
	"testing"

	"github.com/eynopv/lac/internal/assert"
)

func TestReplacePrefixWithEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "Hello")
	defer os.Unsetenv("TEST_ENV_VAR")

	param := Param("${TEST_ENV_VAR}World")
	expected := "HelloWorld"
	result := param.Resolve(nil)
	assert.Equal(t, result, expected)
}

func TestReplacePostfixWithEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "World")
	defer os.Unsetenv("TEST_ENV_VAR")

	param := Param("Hello${TEST_ENV_VAR}")
	expected := "HelloWorld"
	result := param.Resolve(nil)

	assert.Equal(t, result, expected)
}

func TestReplaceMultipleWithEnv(t *testing.T) {
	os.Setenv("TEST_VAR_1", "Hello")
	defer os.Unsetenv("TEST_VAR_1")

	os.Setenv("TEST_VAR_2", "World")
	defer os.Unsetenv("TEST_VAR_2")

	param := Param("${TEST_VAR_1}${TEST_VAR_2}")
	expected := "HelloWorld"
	result := param.Resolve(nil)

	assert.Equal(t, result, expected)
}

func TestReplaceMultipleWithReplacements(t *testing.T) {
	param := Param("${replacement_1}, ${replacement_2}!")
	replacements := map[string]string{
		"replacement_1": "Hello",
		"replacement_2": "World",
	}
	expected := "Hello, World!"
	result := param.Resolve(replacements)

	assert.Equal(t, result, expected)
}
