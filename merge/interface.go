package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/catalyst"
)

type Engine interface {
	ForkchoiceUpdatedV1(heads catalyst.ForkchoiceStateV1, PayloadAttributes *catalyst.PayloadAttributesV1) (catalyst.ForkChoiceResponse, error)
	GetPayloadV1(payloadID hexutil.Bytes) (*catalyst.ExecutableDataV1, error)
	ExecutePayloadV1(params catalyst.ExecutableDataV1) (catalyst.ExecutePayloadResponse, error)
	GetHead() (common.Hash, error)
}
