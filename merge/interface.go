package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Status byte

const (
	VALID Status = iota
	INVALID
	KNOWN
)

type ExecutionPayload struct {
	ParentHash    common.Hash
	Coinbase      common.Address
	StateRoot     common.Hash
	ReceiptRoot   common.Hash
	LogsBloom     types.Bloom
	Random        [32]byte
	BlockNumber   uint64
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	BaseFeePerGas [32]byte
	BlockHash     common.Hash
	Transactions  []types.Transaction
}

type Engine interface {
	PreparePayload(parentHash common.Hash, timestamp uint64, random [32]byte, feeRecipient common.Address, payloadID uint64)
	GetPayload(payloadID uint64) (ExecutionPayload, error)
	ExecutePayload(payload ExecutionPayload) (common.Hash, Status)
	ConsensusValidated(blockHash common.Hash) Status
	ForkchoiceUpdated(headBlockHash, finalizedBlockHash, confirmedBlockHash common.Hash) error
}
