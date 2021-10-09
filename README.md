# Merge-Fuzz

Merge-fuzz is a suite of small tools to fuzz the engine API of Consensus Layer nodes for Ethereum.
It contains several fuzzers that produce interesting inputs to the API and try to crash it.
In the future merge-fuzz will be extended to allow for differential fuzzing of different API implementations
against each other.

## Prerequisite
- Have golang version 1.16 installed (there's an issue with 1.17 and go-fuzz right now)
- Have go-fuzz and go-fuzz-build installed, see github.com/dvyukov/go-fuzz

## Instructions
- Clone the repository `git clone git@github.com:MariusVanDerWijden/merge-fuzz.git`
- Build the fuzzer `cd merge-fuzz/fuzz && CGOENABLED=0 go-fuzz-build`
- Start the test target to listen on localhost:8545 
    (e.g. `geth --dev --catalyst --http --http.api="eth,engine" --override.totalterminaldifficulty=0`)
- Start the fuzzer `go-fuzz --func FunctionName --procs 1`
- Starting `go-fuzz` without a valid function name will tell you which functions are available to fuzz

