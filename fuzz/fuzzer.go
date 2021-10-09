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

var engine merge.Engine

func init() {
	engine = merge.GethRPCEngine //merge.NewGethNode()
}

func FuzzPreparePayload(input []byte) int {
	fuzzer := fuzz.NewFromGoFuzz(input)
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
