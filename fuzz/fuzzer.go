package fuzz

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/trie"
	fuzz "github.com/google/gofuzz"
	"github.com/mariusvanderwijden/merge-fuzz/merge"
)

var engineA merge.Engine
var engineB merge.Engine

func init() {
	engineA = merge.GethRPCEngine //merge.NewGethNode()
	engineB = merge.GethRPCEngine
}

func FuzzPreparePayload(input []byte) int {
	return fuzzPreparePayload(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzGetPayload(input []byte) int { return fuzzGetPayload(fuzz.NewFromGoFuzz(input), engineA) }

func FuzzExecutePayload(input []byte) int {
	return fuzzExecutePayload(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzConsensusValidated(input []byte) int {
	return fuzzConsensusValidated(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzForkchoiceUpdated(input []byte) int {
	return fuzzForkchoiceUpdated(fuzz.NewFromGoFuzz(input), engineA)
}

func FuzzRandom(input []byte) int { return fuzzRandom(fuzz.NewFromGoFuzz(input), engineA) }

func fuzzRandom(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var strategy byte
	fuzzer.Fuzz(&strategy)
	switch strategy % 6 {
	case 0:
		return fuzzPreparePayload(fuzzer, engine)
	case 1:
		return fuzzGetPayload(fuzzer, engine)
	case 2:
		return fuzzExecutePayload(fuzzer, engine)
	case 3:
		return fuzzConsensusValidated(fuzzer, engine)
	case 4:
		return fuzzForkchoiceUpdated(fuzzer, engine)
	case 5:
		return fuzzInteraction(fuzzer, engine)
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
	engine.PreparePayload(catalyst.AssembleBlockParams{ParentHash: parentHash, Timestamp: timestamp, Random: random, FeeRecipient: feeRecipient})
	return 0
}

func fuzzGetPayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var payloadID uint64
	fuzzer.Fuzz(&payloadID)
	payload, err := engine.GetPayload(hexutil.Uint64(payloadID))
	if err != nil {
		return 0
	}
	_ = payload
	return 1
}

func fuzzExecutePayload(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	payload := fillExecPayload(fuzzer)
	_, err := engine.ExecutePayload(payload)
	if err != nil {
		return 1
	}
	return 0
}

func fuzzConsensusValidated(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
	var blockhash common.Hash
	fuzzer.Fuzz(&blockhash)
	err := engine.ConsensusValidated(catalyst.ConsensusValidatedParams{BlockHash: blockhash})
	if err != nil {
		return 1
	}
	return 0
}

func fuzzForkchoiceUpdated(fuzzer *fuzz.Fuzzer, engine merge.Engine) int {
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
	var (
		payload  catalyst.ExecutableData
		realHash bool
		basefee  int64
	)
	fuzzer.Fuzz(&payload)
	fuzzer.Fuzz(&realHash)
	if realHash {
		txs := types.Transactions{}
		fuzzer.Fuzz(&basefee)
		baseFeePerGas := big.NewInt(basefee)
		header := &types.Header{
			ParentHash:  payload.ParentHash,
			UncleHash:   types.EmptyUncleHash,
			Coinbase:    payload.Coinbase,
			Root:        payload.StateRoot,
			TxHash:      types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil)),
			ReceiptHash: payload.ReceiptRoot,
			Bloom:       types.BytesToBloom(payload.LogsBloom),
			Difficulty:  common.Big0,
			Number:      big.NewInt(int64(payload.Number)),
			GasLimit:    payload.GasLimit,
			GasUsed:     payload.GasUsed,
			Time:        payload.Timestamp,
			BaseFee:     baseFeePerGas,
			Extra:       payload.ExtraData,
		}
		payload.Transactions = make([][]byte, 0)
		payload.BlockHash = header.Hash()
		payload.BaseFeePerGas = baseFeePerGas
	}
	return payload
}
