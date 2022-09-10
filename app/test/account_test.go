package test

import "testing"

func add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	if add(1, 2) != 3 {
		t.Errorf("add() = %v, want %v", add(1, 2), 3)
	}
}
