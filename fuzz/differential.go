package fuzz

import (
	"bytes"
	"fmt"
	"os"
	"time"

	fuzz "github.com/google/gofuzz"
)

var logFile = "fuzz_log"

func writeStatus(input []byte) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(string(input) + "\n"); err != nil {
		panic(err)
	}
}

func FuzzDifferential(input []byte) int {
	interesting := false
	writeStatus(input)
	// Execute test on first node
	timestamp := uint64(time.Now().Unix())
	fuzzer := fuzz.NewFromGoFuzz(input)
	a := fuzzRandom(fuzzer, engines[0], timestamp)
	headA, errA := engines[0].GetHead()
	for _, engine := range engines {
		// Execute test on all other nodes
		fuzzer = fuzz.NewFromGoFuzz(input)
		b := fuzzRandom(fuzzer, engine, timestamp)
		headB, errB := engine.GetHead()
		if errA != nil || errB != nil {
			panic(fmt.Sprintf("could not retrieve heads, a: %v, b: %v", errA, errB))
		}

		if !bytes.Equal(headA[:], headB[:]) {
			panic(fmt.Sprintf("different heads, a: %v, b: %v", headA, headB))
		}

		if a != b {
			interesting = true
		}
	}
	if interesting {
		return 1
	}
	return 0
}
