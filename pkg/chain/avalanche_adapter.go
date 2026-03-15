package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// AvalancheAdapter implements the Adapter interface for Avalanche C-Chain
type AvalancheAdapter struct {
	chainID ChainID
}

// NormalizeBlock normalizes an Avalanche block, omitting PoW fields
func (a *AvalancheAdapter) NormalizeBlock(rawBlock *types.Block) (*types.Block, error) {
	if rawBlock == nil {
		return nil, fmt.Errorf("raw block cannot be nil")
	}
	
	// Validate required fields
	if err := validateBlockHash(rawBlock.Hash); err != nil {
		return nil, err
	}
	if err := validateBlockNumber(rawBlock.Number); err != nil {
		return nil, err
	}
	
	normalized := *rawBlock
	
	// Avalanche C-Chain doesn't have PoW fields
	normalized.Difficulty = ""
	normalized.TotalDifficulty = ""
	normalized.MixHash = ""
	normalized.Nonce = ""
	normalized.Sha3Uncles = ""
	normalized.Uncles = nil
	
	return &normalized, nil
}

// GetFinalityDepth returns the number of blocks for finality on Avalanche
func (a *AvalancheAdapter) GetFinalityDepth() int64 {
	// Avalanche has near-instant finality
	return 1
}

// CalculateBlockReward calculates the block reward for Avalanche
func (a *AvalancheAdapter) CalculateBlockReward(block *types.Block) (*big.Int, error) {
	// Avalanche C-Chain block rewards handled by consensus layer
	return big.NewInt(0), nil
}

// NormalizeTransaction normalizes an Avalanche transaction
func (a *AvalancheAdapter) NormalizeTransaction(tx *types.Transaction) (*types.Transaction, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction cannot be nil")
	}
	
	// Validate transaction hash
	if err := validateTransactionHash(tx.Hash); err != nil {
		return nil, err
	}
	
	normalized := *tx
	return &normalized, nil
}

// CalculateEffectiveFee calculates the total fee for an Avalanche transaction
func (a *AvalancheAdapter) CalculateEffectiveFee(tx *types.Transaction, receipt *types.TransactionReceipt) (*big.Int, error) {
	if receipt == nil {
		return nil, fmt.Errorf("receipt cannot be nil for fee calculation")
	}
	
	gasUsed := new(big.Int)
	if _, ok := gasUsed.SetString(receipt.GasUsed, 0); !ok {
		return nil, fmt.Errorf("invalid gasUsed: %s", receipt.GasUsed)
	}
	
	effectiveGasPrice := new(big.Int)
	if tx.EffectiveGasPrice != "" {
		if _, ok := effectiveGasPrice.SetString(tx.EffectiveGasPrice, 0); !ok {
			return nil, fmt.Errorf("invalid effectiveGasPrice: %s", tx.EffectiveGasPrice)
		}
	} else {
		if _, ok := effectiveGasPrice.SetString(tx.GasPrice, 0); !ok {
			return nil, fmt.Errorf("invalid gasPrice: %s", tx.GasPrice)
		}
	}
	
	// Standard EVM fee calculation
	totalFee := new(big.Int).Mul(effectiveGasPrice, gasUsed)
	return totalFee, nil
}

// GetTransactionType identifies the Avalanche transaction type
func (a *AvalancheAdapter) GetTransactionType(tx *types.Transaction) TransactionType {
	return parseTransactionType(tx.Type)
}

// GetL1DataFee returns the L1 data fee (always 0 for Avalanche L1)
func (a *AvalancheAdapter) GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	// Avalanche C-Chain is L1, no L1 data fee
	return big.NewInt(0), nil
}

// IsDepositTransaction checks if a transaction is a deposit (always false for Avalanche L1)
func (a *AvalancheAdapter) IsDepositTransaction(tx *types.Transaction) bool {
	// Avalanche C-Chain has no deposit transactions
	return false
}
