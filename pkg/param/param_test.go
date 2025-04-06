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
	result := param.Resolve(nil, true)
	assert.Equal(t, result, expected)
}

func TestReplacePostfixWithEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "World")
	defer os.Unsetenv("TEST_ENV_VAR")

	param := Param("Hello${TEST_ENV_VAR}")
	expected := "HelloWorld"
	result := param.Resolve(nil, true)

	assert.Equal(t, result, expected)
}

func TestReplaceMultipleWithEnv(t *testing.T) {
	os.Setenv("TEST_VAR_1", "Hello")
	defer os.Unsetenv("TEST_VAR_1")

	os.Setenv("TEST_VAR_2", "World")
	defer os.Unsetenv("TEST_VAR_2")

	param := Param("${TEST_VAR_1}${TEST_VAR_2}")
	expected := "HelloWorld"
	result := param.Resolve(nil, true)

	assert.Equal(t, result, expected)
}

func TestResolveWithoutEnv(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "Hello")
	defer os.Unsetenv("TEST_ENV_VAR")

	param := Param("${TEST_ENV_VAR}World")
	expected := "${TEST_ENV_VAR}World"
	result := param.Resolve(nil, false)

	assert.Equal(t, result, expected)
}

func TestReplaceMultipleWithReplacements(t *testing.T) {
	param := Param("${replacement_1}, ${replacement_2}!")
	replacements := map[string]any{
		"replacement_1": "Hello",
		"replacement_2": "World",
	}
	expected := "Hello, World!"
	result := param.Resolve(replacements, true)

	assert.Equal(t, result, expected)
}

func TestReplaceWithNumericValues(t *testing.T) {
	param := Param("The answer is ${number} and I have ${float_num} dollars")

	replacements := map[string]any{
		"number":    42,
		"float_num": 10.5,
	}

	expected := "The answer is 42 and I have 10.5 dollars"
	result := param.Resolve(replacements, true)

	assert.Equal(t, result, expected)
}

func TestReplaceWithBooleanValues(t *testing.T) {
	param := Param("The statement is ${bool_value}")
	replacements := map[string]any{
		"bool_value": true,
	}
	expected := "The statement is true"
	result := param.Resolve(replacements, true)

	assert.Equal(t, result, expected)
}

func TestReplaceWithNullValues(t *testing.T) {
	param := Param("Is it ${null_value}")
	replacements := map[string]any{
		"null_value": nil,
	}
	expected := "Is it null"
	result := param.Resolve(replacements, true)

	assert.Equal(t, result, expected)
}

func TestReplaceWithQuotedString(t *testing.T) {
	param := Param(`"${value}"`)
	replacements := map[string]any{
		"value": "Hello, World",
	}
	expected := `"Hello, World"`
	result := param.Resolve(replacements, true)

	assert.Equal(t, result, expected)
}

func TestReplaceWithQuotesStringFromEnv(t *testing.T) {
	param := Param(`"${TEST_VAR}"`)

	os.Setenv("TEST_VAR", "Hello, World")
	defer os.Unsetenv("TEST_VAR")

	expected := `"Hello, World"`
	result := param.Resolve(nil, true)

	assert.Equal(t, result, expected)
}
