package fuzz

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/MariusVanDerWijden/FuzzyVM/filler"
	txfuzz "github.com/MariusVanDerWijden/tx-fuzz"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

var engines []merge.Engine

type Config struct {
	Genesis string
	Nodes   []string
	Jwts    []string
}

var once sync.Once

func init() {
	once.Do(func() {
		tomlData, err := ioutil.ReadFile("config.toml")
		if err != nil {
			panic(err)
		}
		var conf Config
		_, err = toml.Decode(string(tomlData), &conf)
		if err != nil {
			panic(err)
		}

		engines = make([]merge.Engine, 0, len(conf.Nodes)+1)
		// TODO remove once differential fuzzing again
		n, _ := merge.NewRPCNode("", "", func() {})
		engines = append(engines, n)
		//engines = append(engines, merge.StartGethNode("genesis.json"))
		for i := range conf.Nodes {
			node, err := merge.NewRPCNode(conf.Nodes[i], conf.Jwts[i], func() {})
			if err != nil {
				panic(err)
			}
			engines = append(engines, node)
		}
	})
}

func FuzzPreparePayload(input []byte) int {
	return fuzzPreparePayload(fuzz.NewFromGoFuzz(input), engines[1])
}

func FuzzGetPayload(input []byte) int { return fuzzGetPayload(fuzz.NewFromGoFuzz(input), engines[1]) }

func FuzzExecutePayload(input []byte) int {
	return fuzzExecutePayload(fuzz.NewFromGoFuzz(input), engines[1])
}

func FuzzForkchoiceUpdated(input []byte) int {
	return fuzzForkchoiceUpdated(fuzz.NewFromGoFuzz(input), engines[1])
}

func FuzzRandom(input []byte) int {
	return fuzzRandom(fuzz.NewFromGoFuzz(input), engines[1], uint64(time.Now().Unix()))
}

func fuzzRandom(fuzzer *fuzz.Fuzzer, engine merge.Engine, timestamp uint64) int {
	var strategy byte
	fuzzer.Fuzz(&strategy)
	switch strategy % 5 {
	case 0:
		return fuzzPreparePayload(fuzzer, engine)
	case 1:
		return fuzzGetPayload(fuzzer, engine)
	case 2:
		return fuzzExecutePayload(fuzzer, engine)
	case 3:
		return fuzzForkchoiceUpdated(fuzzer, engine)
	case 4:
		return fuzzInteraction(fuzzer, engine, timestamp)
	case 5:
		return fuzzSetHead(engine)
	default:
		panic("asdf")
	}
}

func fuzzPreparePayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var (
		parentHash   common.Hash
		timestamp    uint64
		random       [32]byte
		feeRecipient common.Address
		payloadID    uint64
		realHash     bool
	)
	fuzzer.Fuzz(&parentHash)
	fuzzer.Fuzz(&realHash)
	if realHash {
		parentHash, _ = engine.GetHead()
	}
	fuzzer.Fuzz(&timestamp)
	fuzzer.Fuzz(&random)
	fuzzer.Fuzz(&feeRecipient)
	fuzzer.Fuzz(&payloadID)
	heads := beacon.ForkchoiceStateV1{HeadBlockHash: parentHash, SafeBlockHash: parentHash, FinalizedBlockHash: parentHash}
	withdrawals := fuzzWithdrawals(fuzzer)
	attributes := beacon.PayloadAttributes{Timestamp: timestamp, Random: random, SuggestedFeeRecipient: feeRecipient, Withdrawals: withdrawals}
	engine.ForkchoiceUpdatedV2(heads, &attributes)
	return 0
}

func fuzzWithdrawals(fuzzer *fuzz.Fuzzer) types.Withdrawals {
	var cnt byte
	fuzzer.Fuzz(&cnt)
	out := make([]*types.Withdrawal, 0)
	for i := 0; i < int(cnt); i++ {
		var (
			index      uint64
			validator  uint64
			receipient common.Address
			amount     [32]byte
		)
		fuzzer.Fuzz(&index)
		fuzzer.Fuzz(&validator)
		fuzzer.Fuzz(&receipient)
		fuzzer.Fuzz(&amount)

		withdrawal := types.Withdrawal{Index: index, Validator: validator, Address: receipient, Amount: new(big.Int).SetBytes(amount[:])}
		out = append(out, &withdrawal)
	}
	return types.Withdrawals(out)
}

func fuzzGetPayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var payloadID beacon.PayloadID
	fuzzer.Fuzz(&payloadID)
	payload, err := engine.GetPayloadV2(payloadID)
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func fuzzExecutePayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	payload := fillExecPayload(fuzzer)
	_, err := engine.NewPayloadV2(payload)
	if err != nil {
		return 1
	}
	return 0
}

func fuzzForkchoiceUpdated(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var (
		headBlockHash      common.Hash
		safeBlockHash      common.Hash
		finalizedBlockHash common.Hash
	)
	fuzzer.Fuzz(&headBlockHash)
	fuzzer.Fuzz(&safeBlockHash)
	fuzzer.Fuzz(&finalizedBlockHash)
	_, err := engine.ForkchoiceUpdatedV2(beacon.ForkchoiceStateV1{HeadBlockHash: headBlockHash, SafeBlockHash: safeBlockHash, FinalizedBlockHash: finalizedBlockHash}, nil)
	if err == nil {
		return 1
	}
	return 0
}

func fuzzSetHead(engine merge.Engine) int {
	head, _ := engine.GetHead()
	_, err := engine.ForkchoiceUpdatedV2(beacon.ForkchoiceStateV1{HeadBlockHash: head, SafeBlockHash: head, FinalizedBlockHash: head}, nil)
	if err == nil {
		return 1
	}
	return 0
}

func fillExecPayload(fuzzer *fuzz.Fuzzer) beacon.ExecutableData {
	var (
		payload    beacon.ExecutableData
		realHash   bool
		basefee    int64
		fillerData = make([]byte, 128)
		txLen      byte
		txs        []*types.Transaction
	)
	fuzzer.Fuzz(&payload)
	fuzzer.Fuzz(&realHash)
	if realHash {
		fuzzer.Fuzz(&fillerData)
		fuzzer.Fuzz(&txLen)
		f := filler.NewFiller(fillerData)
		node := engines[1].(*merge.RPCnode)
		for i := 0; i < int(txLen); i++ {
			tx, err := txfuzz.RandomValidTx(node.Node, f, payload.FeeRecipient, 0, big.NewInt(0), nil, true)
			if err != nil {
				fmt.Println(err)
				continue
			}
			txs = append(txs, tx)
		}
		fuzzer.Fuzz(&basefee)
		baseFeePerGas := big.NewInt(basefee)
		header := &types.Header{
			ParentHash:  payload.ParentHash,
			UncleHash:   types.EmptyUncleHash,
			Coinbase:    payload.FeeRecipient,
			Root:        payload.StateRoot,
			TxHash:      types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil)),
			ReceiptHash: payload.ReceiptsRoot,
			Bloom:       types.BytesToBloom(payload.LogsBloom),
			Difficulty:  common.Big0,
			Number:      big.NewInt(int64(payload.Number)),
			GasLimit:    payload.GasLimit,
			GasUsed:     payload.GasUsed,
			Time:        payload.Timestamp,
			BaseFee:     baseFeePerGas,
			Extra:       payload.ExtraData,
		}
		payload.Transactions = encodeTransactions(txs)
		payload.BlockHash = header.Hash()
		payload.BaseFeePerGas = baseFeePerGas
	}
	return payload
}

func encodeTransactions(txs []*types.Transaction) [][]byte {
	var enc = make([][]byte, len(txs))
	for i, tx := range txs {
		enc[i], _ = tx.MarshalBinary()
	}
	return enc
}
