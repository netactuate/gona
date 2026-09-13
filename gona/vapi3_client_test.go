package gona

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsV3NotFoundUnwrapsWrappedErrors(t *testing.T) {
	err := fmt.Errorf("get VPC %d: %w", 1, &V3NotFoundError{StatusCode: 404})
	if !IsV3NotFound(err) {
		t.Fatal("expected wrapped V3NotFoundError to be recognized")
	}

	if IsV3NotFound(errors.New("some other error")) {
		t.Fatal("expected unrelated error not to be recognized")
	}
}
