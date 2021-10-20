package fuzz

import (
	"bytes"
	"fmt"

	fuzz "github.com/google/gofuzz"
)

func FuzzDifferential(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	a := fuzzRandom(fuzzer, engineA)
	fuzzer = fuzz.NewFromGoFuzz(input)
	b := fuzzRandom(fuzzer, engineB)
	headA, errA := engineA.GetHead()
	headB, errB := engineB.GetHead()
	if errA != nil || errB != nil {
		panic(fmt.Sprintf("could not retrieve heads, a: %v, b: %v", errA, errB))
	}

	if !bytes.Equal(headA[:], headB[:]) {
		panic(fmt.Sprintf("different heads, a: %v, b: %v", headA, headB))
	}

	if a != b {
		return 1
	}
	return 0
}
