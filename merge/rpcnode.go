package merge

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/rpc"
)

type RPCnode struct {
	node *rpc.Client
}

func NewRPCNode(url string, startNode func()) (*RPCnode, error) {
	startNode()
	node, err := rpc.Dial(url)
	return &RPCnode{node: node}, err
}

func (n *RPCnode) PreparePayload(params catalyst.AssembleBlockParams) (*catalyst.PayloadResponse, error) {
	var res catalyst.PayloadResponse
	err := n.node.Call(&res, "engine_preparePayload", params)
	return &res, err
}

func (n *RPCnode) GetPayload(PayloadID hexutil.Uint64) (*catalyst.ExecutableData, error) {
	var res catalyst.ExecutableData
	err := n.node.Call(&res, "engine_getPayload", PayloadID)
	return &res, err
}

func (n *RPCnode) ExecutePayload(params catalyst.ExecutableData) (catalyst.GenericStringResponse, error) {
	var res catalyst.GenericStringResponse
	err := n.node.Call(&res, "engine_executePayload", params)
	return res, err
}

func (n *RPCnode) ConsensusValidated(params catalyst.ConsensusValidatedParams) error {
	return n.node.Call(nil, "engine_consensusValidated", params)
}

func (n *RPCnode) ForkchoiceUpdated(params catalyst.ForkChoiceParams) error {
	return n.node.Call(nil, "engine_forkchoiceUpdated", params)
}

func (n *RPCnode) GetHead() (common.Hash, error) {
	var head *types.Header
	err := n.node.Call(&head, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		return common.Hash{}, err
	}
	return head.Hash(), nil
}
