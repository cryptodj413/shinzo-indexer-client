package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// OptimismAdapter implements the Adapter interface for Optimism chains
type OptimismAdapter struct {
	chainID ChainID
}

// NormalizeBlock normalizes an Optimism block, omitting PoW fields and including L2 fields
func (a *OptimismAdapter) NormalizeBlock(rawBlock *types.Block) (*types.Block, error) {
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
	
	// Optimism blocks don't have PoW fields
	normalized.Difficulty = ""
	normalized.TotalDifficulty = ""
	normalized.MixHash = ""
	normalized.Nonce = ""
	normalized.Sha3Uncles = ""
	normalized.Uncles = nil
	
	// L2-specific fields are preserved if present
	// (l1BlockNumber)
	
	return &normalized, nil
}

// GetFinalityDepth returns the number of blocks for finality on Optimism
func (a *OptimismAdapter) GetFinalityDepth() int64 {
	// Optimism has instant soft finality on L2
	return 1
}

// CalculateBlockReward calculates the block reward for Optimism
func (a *OptimismAdapter) CalculateBlockReward(block *types.Block) (*big.Int, error) {
	// Optimism doesn't have block rewards (L2 rollup)
	return big.NewInt(0), nil
}

// NormalizeTransaction normalizes an Optimism transaction
func (a *OptimismAdapter) NormalizeTransaction(tx *types.Transaction) (*types.Transaction, error) {
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

// CalculateEffectiveFee calculates the total fee for an Optimism transaction including L1 data fee
func (a *OptimismAdapter) CalculateEffectiveFee(tx *types.Transaction, receipt *types.TransactionReceipt) (*big.Int, error) {
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
	
	// L2 execution fee
	l2Fee := new(big.Int).Mul(effectiveGasPrice, gasUsed)
	
	// Add L1 data fee if available
	if tx.L1DataFee != "" {
		l1DataFee := new(big.Int)
		if _, ok := l1DataFee.SetString(tx.L1DataFee, 0); !ok {
			return nil, fmt.Errorf("invalid l1DataFee: %s", tx.L1DataFee)
		}
		l2Fee.Add(l2Fee, l1DataFee)
	}
	
	return l2Fee, nil
}

// GetTransactionType identifies the Optimism transaction type
func (a *OptimismAdapter) GetTransactionType(tx *types.Transaction) TransactionType {
	txType := parseTransactionType(tx.Type)
	
	// Optimism supports standard types plus deposit transactions
	switch txType {
	case TxTypeLegacy, TxTypeAccessList, TxTypeDynamicFee, TxTypeOptimismDeposit:
		return txType
	default:
		return TxTypeLegacy
	}
}

// GetL1DataFee returns the L1 data fee for an Optimism transaction
func (a *OptimismAdapter) GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	// If already populated in transaction, return it
	if tx.L1DataFee != "" {
		l1DataFee := new(big.Int)
		if _, ok := l1DataFee.SetString(tx.L1DataFee, 0); !ok {
			return nil, fmt.Errorf("invalid l1DataFee: %s", tx.L1DataFee)
		}
		return l1DataFee, nil
	}
	
	// Otherwise, would need to call Optimism-specific RPC method
	// For now, return 0 and let RPC layer populate this
	return big.NewInt(0), nil
}

// IsDepositTransaction checks if a transaction is a deposit transaction (Type 0x7E)
func (a *OptimismAdapter) IsDepositTransaction(tx *types.Transaction) bool {
	return a.GetTransactionType(tx) == TxTypeOptimismDeposit
}
