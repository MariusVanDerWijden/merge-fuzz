package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type RPCnode struct {
	Node *rpc.Client
}

func NewRPCNode(url string, startNode func()) (*RPCnode, error) {
	startNode()
	node, err := rpc.Dial(url)
	return &RPCnode{Node: node}, err
}

func (n *RPCnode) ForkchoiceUpdatedV1(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributesV1) (beacon.ForkChoiceResponse, error) {
	var res beacon.ForkChoiceResponse
	err := n.Node.Call(&res, "engine_forkchoiceUpdatedV1", heads, PayloadAttributes)
	return res, err
}

func (n *RPCnode) NewPayloadV1(params beacon.ExecutableDataV1) (beacon.PayloadStatusV1, error) {
	var res beacon.PayloadStatusV1
	err := n.Node.Call(&res, "engine_newPayloadV1", params)
	return res, err
}

func (n *RPCnode) GetPayloadV1(payloadID beacon.PayloadID) (*beacon.ExecutableDataV1, error) {
	var res beacon.ExecutableDataV1
	err := n.Node.Call(&res, "engine_getPayloadV1", payloadID)
	return &res, err
}

func (n *RPCnode) GetHead() (common.Hash, error) {
	var head *types.Header
	err := n.Node.Call(&head, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		return common.Hash{}, err
	}
	return head.Hash(), nil
}
