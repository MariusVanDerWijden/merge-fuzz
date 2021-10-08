package fuzz

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

var engine merge.Engine

func init() {
	engine = merge.GethRPCEngine //merge.NewGethNode()
}

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
	engine.PreparePayload(catalyst.AssembleBlockParams{ParentHash: parentHash, Timestamp: timestamp, Random: random, FeeRecipient: feeRecipient})
	return 0
}

func FuzzGetPayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var payloadID uint64
	fuzzer.Fuzz(&payloadID)
	payload, err := engine.GetPayload(hexutil.Uint64(payloadID))
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func FuzzExecutePayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	payload := fillExecPayload(fuzzer)
	_, err := engine.ExecutePayload(payload)
	if err != nil {
		return 1
	}
	return 0
}

func FuzzConsensusValidated(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var blockhash common.Hash
	fuzzer.Fuzz(&blockhash)
	err := engine.ConsensusValidated(catalyst.ConsensusValidatedParams{BlockHash: blockhash})
	if err != nil {
		return 1
	}
	return 0
}

func FuzzForkchoiceUpdated(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var (
		headBlockHash      common.Hash
		finalizedBlockHash common.Hash
	)
	fuzzer.Fuzz(&headBlockHash)
	fuzzer.Fuzz(&finalizedBlockHash)
	err := engine.ForkchoiceUpdated(catalyst.ForkChoiceParams{HeadBlockHash: headBlockHash, FinalizedBlockHash: finalizedBlockHash})
	if err == nil {
		return 1
	}
	return 0
}

func fillExecPayload(fuzzer *fuzz.Fuzzer) catalyst.ExecutableData {
	var payload catalyst.ExecutableData
	fuzzer.Fuzz(&payload)
	return payload
}
