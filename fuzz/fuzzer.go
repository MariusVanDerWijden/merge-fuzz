package fuzz

import (
	"fmt"
	"math/big"
	"time"

	"github.com/MariusVanDerWijden/FuzzyVM/filler"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/trie"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
	txfuzz "github.com/mariusvanderwijden/tx-fuzz"
)

var engineA merge.Engine
var engineB merge.Engine

func init() {
	engineA, _ = merge.NewRPCNode("http://127.0.0.1:8545", func() {}) //merge.NewGethNode()
	engineB, _ = merge.NewRPCNode("http://127.0.0.1:8546", func() {})
}

func FuzzPreparePayload(input []byte) int {
	return fuzzPreparePayload(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzGetPayload(input []byte) int { return fuzzGetPayload(fuzz.NewFromGoFuzz(input), engineA) }

func FuzzExecutePayload(input []byte) int {
	return fuzzExecutePayload(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzForkchoiceUpdated(input []byte) int {
	return fuzzForkchoiceUpdated(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzRandom(input []byte) int {
	return fuzzRandom(fuzz.NewFromGoFuzz(input), engineA, uint64(time.Now().Unix()))
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
	heads := catalyst.ForkchoiceStateV1{HeadBlockHash: parentHash, SafeBlockHash: parentHash, FinalizedBlockHash: parentHash}
	attributes := catalyst.PayloadAttributesV1{Timestamp: timestamp, Random: random, SuggestedFeeRecipient: feeRecipient}
	engine.ForkchoiceUpdatedV1(heads, &attributes)
	return 0
}

func fuzzGetPayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var payloadID hexutil.Bytes
	fuzzer.Fuzz(&payloadID)
	payload, err := engine.GetPayloadV1(payloadID)
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func fuzzExecutePayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	payload := fillExecPayload(fuzzer)
	_, err := engine.ExecutePayloadV1(payload)
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
	_, err := engine.ForkchoiceUpdatedV1(catalyst.ForkchoiceStateV1{HeadBlockHash: headBlockHash, SafeBlockHash: safeBlockHash, FinalizedBlockHash: finalizedBlockHash}, nil)
	if err == nil {
		return 1
	}
	return 0
}

func fuzzSetHead(engine merge.Engine) int {
	head, _ := engine.GetHead()
	_, err := engine.ForkchoiceUpdatedV1(catalyst.ForkchoiceStateV1{HeadBlockHash: head, SafeBlockHash: head, FinalizedBlockHash: head}, nil)
	if err == nil {
		return 1
	}
	return 0
}

func fillExecPayload(fuzzer *fuzz.Fuzzer) catalyst.ExecutableDataV1 {
	var (
		payload    catalyst.ExecutableDataV1
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
		node := engineA.(*merge.RPCnode)
		for i := 0; i < int(txLen); i++ {
			tx, err := txfuzz.RandomValidTx(node.Node, f, payload.FeeRecipient, 0, big.NewInt(0), nil)
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
