package internal

import "testing"

func TestContext(t *testing.T) {
	context := GetContext()
	headers := map[string]string{
		"Hello": "World",
	}
	context.Config = &Config{Headers: headers}

	contextOther := GetContext()
	if _, ok := contextOther.Config.Headers["Hello"]; ok == false {
		t.Fatalf("Expected context to have config")
	}
}
