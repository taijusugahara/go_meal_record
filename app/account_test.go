package main

import "testing"

func add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	if add(1, 2) != 4 {
		t.Errorf("add() = %v, want %v", add(1, 2), 3)
	}
}
