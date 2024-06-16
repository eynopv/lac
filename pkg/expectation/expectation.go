package expectation

import (
	"fmt"
	"time"

	"github.com/eynopv/lac/pkg/result"
)

type Expectation struct {
	Status       int           `json:"status"`
	TimeLessThan time.Duration `json:"timeLessThan"`
}

func (e Expectation) Check(r *result.Result) error {
	if e.Status != 0 && e.Status != r.StatusCode {
		return fmt.Errorf("Expected status %v but got %v", e.Status, r.StatusCode)
	}

	if e.TimeLessThan != 0 && e.TimeLessThan*time.Millisecond < r.ElapsedTime {
		return fmt.Errorf(
			"Expected duration less than %v but got %v", e.TimeLessThan, r.ElapsedTime,
		)
	}

	return nil
}
