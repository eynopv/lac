package expectation

import (
	"fmt"
	"testing"
	"time"

	"github.com/eynopv/lac/internal/result"
)

func TestSuccessfullExpectation(t *testing.T) {
	expect := Expectation{Status: 200, TimeLessThan: 100}
	result := result.Result{StatusCode: 200, ElapsedTime: 100 * time.Millisecond}

	err := expect.Check(&result)
	if err != nil {
		t.Fatalf("Expected error to be nil")
	}
}

func TestFailedStatusExpectation(t *testing.T) {
	expect := Expectation{Status: 200, TimeLessThan: 100}
	result := result.Result{StatusCode: 400, ElapsedTime: 99 * time.Millisecond}

	err := expect.Check(&result)
	if err == nil {
		t.Fatalf("Expected failed expectation by status")
	}

	expectedMessage := fmt.Sprintf("Expected status %v but got %v", expect.Status, result.StatusCode)

	if err.Error() != expectedMessage {
		t.Fatalf("\nExpected: " + expectedMessage + "\nReceived: " + err.Error() + "\n")
	}
}

func TestFailedTimeExpectation(t *testing.T) {
	expect := Expectation{Status: 200, TimeLessThan: 100}
	result := result.Result{StatusCode: 200, ElapsedTime: 110 * time.Millisecond}

	err := expect.Check(&result)
	if err == nil {
		t.Fatalf("Expected failed expectation by time")
	}

	expectedMessage := fmt.Sprintf(
		"Expected duration less than %v but got %v",
		expect.TimeLessThan,
		result.ElapsedTime,
	)

	if err.Error() != expectedMessage {
		t.Fatalf("\nExpected: " + expectedMessage + "\nReceived: " + err.Error() + "\n")
	}
}
