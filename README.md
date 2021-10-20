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

## Example for differential fuzzing
In order to differential fuzz two targets against each other, they need to be started on `localhost:8545` and `localhost:8546`.
They also need to be initialized with the same genesis block
An example of differential fuzzing two geth instances against each other:
```
rm -rf fuzz1/ && rm -rf fuzz2/

~/go/src/github.com/ethereum/go-ethereum/build/bin/geth init genesis.json  --datadir "fuzz1"
~/go/src/github.com/ethereum/go-ethereum/build/bin/geth init genesis.json  --datadir "fuzz2"


~/go/src/github.com/ethereum/go-ethereum/build/bin/geth --datadir "fuzz1" --catalyst --http --http.api="eth,engine" --http.port 8545 --override.totalterminaldifficulty=0
~/go/src/github.com/ethereum/go-ethereum/build/bin/geth --datadir "fuzz2" --catalyst --http --http.api="eth,engine" --http.port 8546 -port 30304 --override.totalterminaldifficulty=0
```