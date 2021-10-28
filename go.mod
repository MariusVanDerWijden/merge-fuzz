module github.com/mariusvanderwijden/merge-fuzz

go 1.16

require (
	github.com/MariusVanDerWijden/FuzzyVM v0.0.0-20210904205340-da82a0d3e27a
	github.com/ethereum/go-ethereum v1.10.11
	github.com/google/gofuzz v1.2.0
	github.com/mariusvanderwijden/tx-fuzz v0.0.0-20211025152518-e9a05306f573
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
)

replace github.com/ethereum/go-ethereum => /home/matematik/go/src/github.com/ethereum/go-ethereum
