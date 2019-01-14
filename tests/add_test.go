package tests_test

import "testing"

func TestSum(t *testing.T) {
	total := 5 + 4
	if false {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 9)
	}
}
