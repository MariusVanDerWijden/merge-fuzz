package fuzz

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

func FuzzInteraction(input []byte) int {
	return fuzzInteraction(fuzz.NewFromGoFuzz(input), engineA, uint64(time.Now().Unix()))
}

func fuzzInteraction(fuzzer *fuzz.Fuzzer, engine merge.Engine, timestamp uint64) int {
	var (
		random        [32]byte
		feeRecipient  common.Address
		realTimestamp byte
	)
	parentHash, err := engine.GetHead()
	if err != nil {
		panic(err)
	}
	fuzzer.Fuzz(&random)
	fuzzer.Fuzz(&feeRecipient)
	fuzzer.Fuzz(&realTimestamp)
	if realTimestamp > 30 {
		timestamp += 12
	} else if realTimestamp > 60 {
		timestamp -= 12
	}
	response, err := engine.ForkchoiceUpdatedV1(beacon.ForkchoiceStateV1{HeadBlockHash: parentHash, SafeBlockHash: parentHash, FinalizedBlockHash: parentHash}, &beacon.PayloadAttributesV1{Timestamp: timestamp, Random: random, SuggestedFeeRecipient: feeRecipient})
	if err != nil {
		return 0
	}
	payload, err := engine.GetPayloadV1(*response.PayloadID)
	if err != nil {
		return 0
	}
	resp1, err := engine.NewPayloadV1(*payload)
	if err != nil {
		panic(err)
	}
	resp2, err := engine.NewPayloadV1(*payload)
	if err != nil {
		panic(err)
	}
	if resp1.Status != resp2.Status {
		panic(fmt.Sprintf("invalid status %v %v", resp1, resp2))
	}
	response, err = engine.ForkchoiceUpdatedV1(beacon.ForkchoiceStateV1{HeadBlockHash: payload.BlockHash, SafeBlockHash: payload.BlockHash, FinalizedBlockHash: payload.BlockHash}, nil)
	if err != nil {
		return 0
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
