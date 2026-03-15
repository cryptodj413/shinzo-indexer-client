package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// ArbitrumAdapter implements the Adapter interface for Arbitrum chains
type ArbitrumAdapter struct {
	chainID ChainID
}

// NormalizeBlock normalizes an Arbitrum block, omitting PoW fields and including L2 fields
func (a *ArbitrumAdapter) NormalizeBlock(rawBlock *types.Block) (*types.Block, error) {
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
	
	// Arbitrum blocks don't have PoW fields
	normalized.Difficulty = ""
	normalized.TotalDifficulty = ""
	normalized.MixHash = ""
	normalized.Nonce = ""
	normalized.Sha3Uncles = ""
	normalized.Uncles = nil
	
	// L2-specific fields are preserved if present
	// (l1BlockNumber, sendCount, sendRoot)
	
	return &normalized, nil
}

// GetFinalityDepth returns the number of blocks for finality on Arbitrum
func (a *ArbitrumAdapter) GetFinalityDepth() int64 {
	// Arbitrum has instant soft finality
	return 1
}

// CalculateBlockReward calculates the block reward for Arbitrum
func (a *ArbitrumAdapter) CalculateBlockReward(block *types.Block) (*big.Int, error) {
	// Arbitrum doesn't have block rewards (L2 rollup)
	return big.NewInt(0), nil
}

// NormalizeTransaction normalizes an Arbitrum transaction
func (a *ArbitrumAdapter) NormalizeTransaction(tx *types.Transaction) (*types.Transaction, error) {
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

// CalculateEffectiveFee calculates the total fee for an Arbitrum transaction including L1 data fee
func (a *ArbitrumAdapter) CalculateEffectiveFee(tx *types.Transaction, receipt *types.TransactionReceipt) (*big.Int, error) {
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

// GetTransactionType identifies the Arbitrum transaction type
func (a *ArbitrumAdapter) GetTransactionType(tx *types.Transaction) TransactionType {
	txType := parseTransactionType(tx.Type)
	
	// Arbitrum supports standard types plus retryable tickets
	switch txType {
	case TxTypeLegacy, TxTypeAccessList, TxTypeDynamicFee, TxTypeArbitrumRetryable:
		return txType
	default:
		return TxTypeLegacy
	}
}

// GetL1DataFee returns the L1 data fee for an Arbitrum transaction
func (a *ArbitrumAdapter) GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	// If already populated in transaction, return it
	if tx.L1DataFee != "" {
		l1DataFee := new(big.Int)
		if _, ok := l1DataFee.SetString(tx.L1DataFee, 0); !ok {
			return nil, fmt.Errorf("invalid l1DataFee: %s", tx.L1DataFee)
		}
		return l1DataFee, nil
	}
	
	// Otherwise, would need to call Arbitrum-specific RPC method
	// For now, return 0 and let RPC layer populate this
	return big.NewInt(0), nil
}

// IsDepositTransaction checks if a transaction is a retryable ticket (deposit)
func (a *ArbitrumAdapter) IsDepositTransaction(tx *types.Transaction) bool {
	return a.GetTransactionType(tx) == TxTypeArbitrumRetryable
}
