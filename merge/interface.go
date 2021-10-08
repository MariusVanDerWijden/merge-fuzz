package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/catalyst"
)

type Engine interface {
	PreparePayload(params catalyst.AssembleBlockParams) (*catalyst.PayloadResponse, error)
	GetPayload(PayloadID hexutil.Uint64) (*catalyst.ExecutableData, error)
	ExecutePayload(params catalyst.ExecutableData) (catalyst.GenericStringResponse, error)
	ConsensusValidated(params catalyst.ConsensusValidatedParams) error
	ForkchoiceUpdated(params catalyst.ForkChoiceParams) error
	GetHead() common.Hash
}
