package fuzz

import "testing"

func TestFuzz(t *testing.T) {
	input := "1234adsfasdf"
	FuzzInteraction([]byte(input))
}
