package merge

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	jwt "github.com/golang-jwt/jwt/v4"
)

type RPCnode struct {
	Node *rpc.Client
	Jwt  []byte
}

func NewRPCNode(url, secret string, startNode func()) (*RPCnode, error) {
	startNode()
	node, err := rpc.Dial(url)

	return &RPCnode{Node: node, Jwt: common.Hex2Bytes(secret)}, err
}

func issueToken(secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
	})
	ss, _ := token.SignedString(secret)
	return ss
}

func (n *RPCnode) ForkchoiceUpdatedV1(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributes) (beacon.ForkChoiceResponse, error) {
	var res beacon.ForkChoiceResponse
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&res, "engine_forkchoiceUpdatedV1", heads, PayloadAttributes)
	return res, err
}

func (n *RPCnode) ForkchoiceUpdatedV2(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributes) (beacon.ForkChoiceResponse, error) {
	var res beacon.ForkChoiceResponse
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&res, "engine_forkchoiceUpdatedV2", heads, PayloadAttributes)
	return res, err
}

func (n *RPCnode) NewPayloadV1(params beacon.ExecutableData) (beacon.PayloadStatusV1, error) {
	var res beacon.PayloadStatusV1
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&res, "engine_newPayloadV1", params)
	return res, err
}

func (n *RPCnode) GetPayloadV1(payloadID beacon.PayloadID) (*beacon.ExecutableData, error) {
	var res beacon.ExecutableData
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&res, "engine_getPayloadV1", payloadID)
	return &res, err
}

func (n *RPCnode) NewPayloadV2(params beacon.ExecutableData) (beacon.PayloadStatusV1, error) {
	var res beacon.PayloadStatusV1
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&res, "engine_newPayloadV2", params)
	return res, err
}

func (n *RPCnode) GetPayloadV2(payloadID beacon.PayloadID) (*beacon.ExecutableData, error) {
	var res beacon.ExecutableData
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&res, "engine_getPayloadV2", payloadID)
	return &res, err
}

func (n *RPCnode) GetHead() (common.Hash, error) {
	var head *types.Header
	n.Node.SetHeader("Authorization", "Bearer "+issueToken(n.Jwt))
	err := n.Node.Call(&head, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		return common.Hash{}, err
	}
	return head.Hash(), nil
}
