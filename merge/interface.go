package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
)

type Engine interface {
	ForkchoiceUpdatedV1(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributes) (beacon.ForkChoiceResponse, error)
	ForkchoiceUpdatedV2(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributes) (beacon.ForkChoiceResponse, error)
	GetPayloadV1(payloadID beacon.PayloadID) (*beacon.ExecutableData, error)
	NewPayloadV1(params beacon.ExecutableData) (beacon.PayloadStatusV1, error)
	GetPayloadV2(payloadID beacon.PayloadID) (*beacon.ExecutableData, error)
	NewPayloadV2(params beacon.ExecutableData) (beacon.PayloadStatusV1, error)
	GetHead() (common.Hash, error)
}
