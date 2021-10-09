module github.com/mariusvanderwijden/merge-fuzz

go 1.16

require (
	github.com/dvyukov/go-fuzz v0.0.0-20210914135545-4980593459a1 // indirect
	github.com/ethereum/go-ethereum v1.10.8
	github.com/google/gofuzz v1.2.0
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
)

replace github.com/ethereum/go-ethereum => /home/matematik/go/src/github.com/ethereum/go-ethereum
