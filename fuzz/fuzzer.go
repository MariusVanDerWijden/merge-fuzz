package fuzz

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress(testKey.PublicKey)

	testBalance = big.NewInt(2e18)

	engine merge.Engine
)

func init() {
	genesis, blocks := generatePreMergeChain(10)
	_, ethservice, err := startEthService(genesis, blocks)
	if err != nil {
		panic(err)
	}
	engine = catalyst.NewConsensusAPI(ethservice, nil)
}

func generatePreMergeChain(n int) (*core.Genesis, []*types.Block) {
	db := rawdb.NewMemoryDatabase()
	config := params.AllEthashProtocolChanges
	genesis := &core.Genesis{
		Config:    config,
		Alloc:     core.GenesisAlloc{testAddr: {Balance: testBalance}},
		ExtraData: []byte("test genesis"),
		Timestamp: 9000,
		BaseFee:   big.NewInt(params.InitialBaseFee),
	}
	testNonce := uint64(0)
	generate := func(i int, g *core.BlockGen) {
		g.OffsetTime(5)
		g.SetExtra([]byte("test"))
		tx, _ := types.SignTx(types.NewTransaction(testNonce, common.HexToAddress("0x9a9070028361F7AAbeB3f2F2Dc07F82C4a98A02a"), big.NewInt(1), params.TxGas, big.NewInt(params.InitialBaseFee*2), nil), types.LatestSigner(config), testKey)
		g.AddTx(tx)
		testNonce++
	}
	gblock := genesis.ToBlock(db)
	engine := ethash.NewFaker()
	blocks, _ := core.GenerateChain(config, gblock, engine, db, n, generate)
	totalDifficulty := big.NewInt(0)
	for _, b := range blocks {
		totalDifficulty.Add(totalDifficulty, b.Difficulty())
	}
	config.TerminalTotalDifficulty = totalDifficulty
	return genesis, blocks
}

func startEthService(genesis *core.Genesis, blocks []*types.Block) (*node.Node, *eth.Ethereum, error) {

	n, err := node.New(&node.Config{})
	if err != nil {
		return nil, nil, err
	}

	ethcfg := &ethconfig.Config{Genesis: genesis, Ethash: ethash.Config{PowMode: ethash.ModeFake}, TrieTimeout: time.Minute, TrieDirtyCache: 256, TrieCleanCache: 256}
	ethservice, err := eth.New(n, ethcfg)
	if err != nil {
		return nil, nil, err
	}
	if err := n.Start(); err != nil {
		return nil, nil, err
	}
	if _, err := ethservice.BlockChain().InsertChain(blocks); err != nil {
		n.Close()
		return nil, nil, err
	}
	ethservice.SetEtherbase(testAddr)
	ethservice.SetSynced()

	return n, ethservice, nil
}

func FuzzPreparePayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var (
		parentHash   common.Hash
		timestamp    uint64
		random       [32]byte
		feeRecipient common.Address
		payloadID    uint64
	)
	fuzzer.Fuzz(&parentHash)
	fuzzer.Fuzz(&timestamp)
	fuzzer.Fuzz(&random)
	fuzzer.Fuzz(&feeRecipient)
	fuzzer.Fuzz(&payloadID)
	engine.PreparePayload(catalyst.AssembleBlockParams{ParentHash: parentHash, Timestamp: timestamp, Random: random, FeeRecipient: feeRecipient})
	return 0
}

func FuzzGetPayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var payloadID uint64
	fuzzer.Fuzz(&payloadID)
	payload, err := engine.GetPayload(hexutil.Uint64(payloadID))
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func FuzzExecutePayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	payload := fillExecPayload(fuzzer)
	_, err := engine.ExecutePayload(payload)
	if err != nil {
		return 1
	}
	return 0
}

func FuzzConsensusValidated(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var blockhash common.Hash
	fuzzer.Fuzz(&blockhash)
	err := engine.ConsensusValidated(catalyst.ConsensusValidatedParams{BlockHash: blockhash})
	if err != nil {
		return 1
	}
	return 0
}

func FuzzForkchoiceUpdated(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
	var (
		headBlockHash      common.Hash
		finalizedBlockHash common.Hash
	)
	fuzzer.Fuzz(&headBlockHash)
	fuzzer.Fuzz(&finalizedBlockHash)
	err := engine.ForkchoiceUpdated(catalyst.ForkChoiceParams{HeadBlockHash: headBlockHash, FinalizedBlockHash: finalizedBlockHash})
	if err == nil {
		return 1
	}
	return 0
}

func fillExecPayload(fuzzer *fuzz.Fuzzer) catalyst.ExecutableData {
	var payload catalyst.ExecutableData
	fuzzer.Fuzz(&payload)
	return payload
}
