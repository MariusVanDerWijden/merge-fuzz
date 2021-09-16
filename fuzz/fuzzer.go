package fuzz

import (
	"github.com/ethereum/go-ethereum/common"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

var engine merge.Engine

func FuzzPreparePayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var (
		parentHash   common.Hash
		timestamp    uint64
		random       [32]byte
		feeRecipient common.Address
		payloadID    uint64
	)
	fuzzer.Fuzz(&parentHash)
	fuzzer.Fuzz(&timestamp)
	fuzzer.Fuzz(&random)
	fuzzer.Fuzz(&feeRecipient)
	fuzzer.Fuzz(&payloadID)
	engine.PreparePayload(parentHash, timestamp, random, feeRecipient, payloadID)
	return 0
}

func FuzzGetPayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var payloadID uint64
	fuzzer.Fuzz(&payloadID)
	payload, err := engine.GetPayload(payloadID)
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func FuzzExecutePayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	payload := fillExecPayload(fuzzer)
	hash, status := engine.ExecutePayload(payload)
	if status == merge.VALID || status == merge.KNOWN {
		return 1
	}
	_ = hash
	return 0
}

func FuzzConsensusValidated(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var blockhash common.Hash
	fuzzer.Fuzz(&blockhash)
	status := engine.ConsensusValidated(blockhash)
	if status == merge.VALID || status == merge.KNOWN {
		return 1
	}
	return 0
}

func FuzzForkchoiceUpdated(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var (
		headBlockHash      common.Hash
		finalizedBlockHash common.Hash
		confirmedBlockHash common.Hash
	)
	fuzzer.Fuzz(&headBlockHash)
	fuzzer.Fuzz(&finalizedBlockHash)
	fuzzer.Fuzz(&confirmedBlockHash)
	err := engine.ForkchoiceUpdated(headBlockHash, finalizedBlockHash, confirmedBlockHash)
	if err == nil {
		return 1
	}
	return 0
}

func fillExecPayload(fuzzer *fuzz.Fuzzer) merge.ExecutionPayload {
	var payload merge.ExecutionPayload
	fuzzer.Fuzz(&payload)
	return payload
}

func FuzzInteraction(input []byte) int {
	return 0
}
