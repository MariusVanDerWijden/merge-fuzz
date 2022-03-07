package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
)

type Engine interface {
	ForkchoiceUpdatedV1(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributesV1) (beacon.ForkChoiceResponse, error)
	GetPayloadV1(payloadID beacon.PayloadID) (*beacon.ExecutableDataV1, error)
	ExecutePayloadV1(params beacon.ExecutableDataV1) (beacon.ExecutePayloadResponse, error)
	GetHead() (common.Hash, error)
}
