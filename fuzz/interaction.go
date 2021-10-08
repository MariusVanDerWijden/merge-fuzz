package fuzz

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	fuzz "github.com/google/gofuzz"
)

func FuzzInteraction(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var (
		timestamp    uint64
		random       [32]byte
		feeRecipient common.Address
	)
	// TODO set valid parent hash
	// fuzzer.Fuzz(&parentHash)
	parentHash, err := engine.GetHead()
	if err != nil {
		panic(err)
	}
	fuzzer.Fuzz(&timestamp)
	fuzzer.Fuzz(&random)
	fuzzer.Fuzz(&feeRecipient)
	payloadID, err := engine.PreparePayload(catalyst.AssembleBlockParams{ParentHash: parentHash, Timestamp: timestamp, Random: random, FeeRecipient: feeRecipient})
	if err != nil {
		return 0
	}
	payload, err := engine.GetPayload(hexutil.Uint64(payloadID.PayloadID))
	if err != nil {
		return 0
	}
	resp1, err := engine.ExecutePayload(*payload)
	if err != nil {
		panic(err)
	}
	resp2, err := engine.ExecutePayload(*payload)
	if err != nil {
		panic(err)
	}
	if resp1.Status != resp2.Status {
		panic(fmt.Sprintf("invalid status %v %v", resp1, resp2))
	}
	err = engine.ConsensusValidated(catalyst.ConsensusValidatedParams{BlockHash: payload.BlockHash, Status: catalyst.VALID.Status})
	if err != nil {
		panic(err)
	}
	err = engine.ForkchoiceUpdated(catalyst.ForkChoiceParams{HeadBlockHash: payload.BlockHash, FinalizedBlockHash: common.Hash{}})
	if err != nil {
		panic(err)
	}
	// check that head is updated
	newHead, err := engine.GetHead()
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(newHead[:], payload.BlockHash[:]) {
		panic(fmt.Errorf("invalid head: got %v want %v", newHead, payload.BlockHash))
	}
	return 0
}
