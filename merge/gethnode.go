package merge

import (
	"encoding/json"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/node"
)

type gethNode struct {
	eth *eth.Ethereum
	api catalyst.ConsensusAPI
}

func StartGethNode(filename string) *gethNode {
	// import genesis
	genesis := new(core.Genesis)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		panic(err)
	}
	// Create node
	n, err := node.New(&node.Config{HTTPPort: 1234, AuthPort: 1235})
	if err != nil {
		panic(err)
	}

	ethcfg := &ethconfig.Config{Genesis: genesis, TrieTimeout: time.Minute, TrieDirtyCache: 256, TrieCleanCache: 256, Miner: miner.Config{GasCeil: 30_000_000}}
	ethservice, err := eth.New(n, ethcfg)
	if err != nil {
		panic(err)
	}
	if err := n.Start(); err != nil {
		panic(err)
	}
	return &gethNode{
		eth: ethservice,
		api: *catalyst.NewConsensusAPI(ethservice),
	}
}

func (g *gethNode) ForkchoiceUpdatedV1(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributes) (beacon.ForkChoiceResponse, error) {
	return g.api.ForkchoiceUpdatedV1(heads, PayloadAttributes)
}

func (g *gethNode) ForkchoiceUpdatedV2(heads beacon.ForkchoiceStateV1, PayloadAttributes *beacon.PayloadAttributes) (beacon.ForkChoiceResponse, error) {
	return g.api.ForkchoiceUpdatedV2(heads, PayloadAttributes)
}

func (g *gethNode) GetPayloadV1(payloadID beacon.PayloadID) (*beacon.ExecutableData, error) {
	return g.api.GetPayloadV1(payloadID)
}

func (g *gethNode) NewPayloadV1(params beacon.ExecutableData) (beacon.PayloadStatusV1, error) {
	return g.api.NewPayloadV1(params)
}

func (g *gethNode) GetPayloadV2(payloadID beacon.PayloadID) (*beacon.ExecutableData, error) {
	return g.api.GetPayloadV2(payloadID)
}

func (g *gethNode) NewPayloadV2(params beacon.ExecutableData) (beacon.PayloadStatusV1, error) {
	return g.api.NewPayloadV2(params)
}

func (g *gethNode) GetHead() (common.Hash, error) {
	return g.eth.BlockChain().CurrentHeader().Hash(), nil
}
