package tests_test

import "testing"

func TestSum(t *testing.T) {
	total := 5 + 4
	if total != 9 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 9)
	}
}
