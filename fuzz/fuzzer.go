package fuzz

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

var engine merge.Engine

func FuzzPreparePayload(input []byte) int {
	var (
		parentHash   common.Hash
		timestamp    uint64
		random       [32]byte
		feeRecipient common.Address
		payload      uint64
	)
	engine.PreparePayload(parentHash, timestamp, random, feeRecipient, payload)
	return 0
}

func FuzzGetPayload(input []byte) int {
	var payloadID uint64
	payload, err := engine.GetPayload(payloadID)
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func FuzzExecutePayload(input []byte) int {
	var payload merge.ExecutionPayload
	hash, status := engine.ExecutePayload(payload)
	if status == merge.VALID || status == merge.KNOWN {
		return 1
	}
	_ = hash
	return 0
}

func FuzzConsensusValidated(input []byte) int {
	var blockhash common.Hash
	status := engine.ConsensusValidated(blockhash)
	if status == merge.VALID || status == merge.KNOWN {
		return 1
	}
	return 0
}

func FuzzForkchoiceUpdated(input []byte) int {
	var (
		headBlockHash      common.Hash
		finalizedBlockHash common.Hash
		confirmedBlockHash common.Hash
	)
	err := engine.ForkchoiceUpdated(headBlockHash, finalizedBlockHash, confirmedBlockHash)
	if err == nil {
		return 1
	}
	return 0
}
