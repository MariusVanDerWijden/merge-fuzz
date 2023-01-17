package fuzz

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestInteraction(t *testing.T) {
	engine := engines[1]
	parentHash, err := engine.GetHead()
	if err != nil {
		panic(err)
	}
	parentHash = common.HexToHash("0xfe950635b1bd2a416ff6283b0bbd30176e1b1125ad06fa729da9f3f4c1c61710")
	random := [32]byte{0xff}
	feeRecipient := common.Address{0xaa}
	timestamp := uint64(time.Now().Second() + 1)

	out := make([]*types.Withdrawal, 0)
	for i := 0; i < int(8); i++ {
		index := uint64(i)
		validator := uint64(i + 0xffff)
		receipient := common.Address{byte(i)}
		amount := [32]byte{byte(i)}

		withdrawal := types.Withdrawal{Index: index, Validator: validator, Address: receipient, Amount: new(big.Int).SetBytes(amount[:])}
		out = append(out, &withdrawal)
	}
	withdrawals := types.Withdrawals(out)
	response, err := engine.ForkchoiceUpdatedV2(beacon.ForkchoiceStateV1{HeadBlockHash: parentHash, SafeBlockHash: parentHash, FinalizedBlockHash: parentHash}, &beacon.PayloadAttributes{Timestamp: timestamp, Random: random, SuggestedFeeRecipient: feeRecipient, Withdrawals: withdrawals})
	if err != nil {
		panic(err)
	}
	payload, err := engine.GetPayloadV2(*response.PayloadID)
	if err != nil {
		panic(err)
	}
	_, err = engine.NewPayloadV2(*payload)
	if err != nil {
		panic(err)
	}
	response, err = engine.ForkchoiceUpdatedV2(beacon.ForkchoiceStateV1{HeadBlockHash: payload.BlockHash, SafeBlockHash: payload.BlockHash, FinalizedBlockHash: payload.BlockHash}, nil)
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
	panic("asdf")
}
