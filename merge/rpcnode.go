package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/catalyst"
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

func (n *RPCnode) ForkchoiceUpdatedV1(heads catalyst.ForkchoiceStateV1, PayloadAttributes *catalyst.PayloadAttributesV1) (catalyst.ForkChoiceResponse, error) {
	var res catalyst.ForkChoiceResponse
	err := n.Node.Call(&res, "engine_forkchoiceUpdatedV1", heads, PayloadAttributes)
	return res, err
}

func (n *RPCnode) ExecutePayloadV1(params catalyst.ExecutableDataV1) (catalyst.ExecutePayloadResponse, error) {
	var res catalyst.ExecutePayloadResponse
	err := n.Node.Call(&res, "engine_executePayloadV1", params)
	return res, err
}

func (n *RPCnode) GetPayloadV1(payloadID hexutil.Bytes) (*catalyst.ExecutableDataV1, error) {
	var res catalyst.ExecutableDataV1
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
