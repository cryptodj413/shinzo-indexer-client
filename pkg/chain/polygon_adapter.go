package chain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// PolygonAdapter implements the Adapter interface for Polygon PoS chains
type PolygonAdapter struct {
	chainID ChainID
}

// NormalizeBlock normalizes a Polygon block
func (a *PolygonAdapter) NormalizeBlock(rawBlock *types.Block) (*types.Block, error) {
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
	
	// Polygon has difficulty field but it's always 0 (sidechain)
	// Keep the field but it won't be meaningful
	
	return &normalized, nil
}

// GetFinalityDepth returns the number of blocks for finality on Polygon
func (a *PolygonAdapter) GetFinalityDepth() int64 {
	// Polygon requires 128 blocks for Ethereum checkpoint finality
	return 128
}

// CalculateBlockReward calculates the block reward for Polygon
func (a *PolygonAdapter) CalculateBlockReward(block *types.Block) (*big.Int, error) {
	// Polygon PoS validators receive rewards
	// This is a simplified implementation
	return big.NewInt(0), nil
}

// NormalizeTransaction normalizes a Polygon transaction
func (a *PolygonAdapter) NormalizeTransaction(tx *types.Transaction) (*types.Transaction, error) {
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

// CalculateEffectiveFee calculates the total fee for a Polygon transaction
func (a *PolygonAdapter) CalculateEffectiveFee(tx *types.Transaction, receipt *types.TransactionReceipt) (*big.Int, error) {
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

// GetTransactionType identifies the Polygon transaction type
func (a *PolygonAdapter) GetTransactionType(tx *types.Transaction) TransactionType {
	return parseTransactionType(tx.Type)
}

// GetL1DataFee returns the L1 data fee (always 0 for Polygon sidechain)
func (a *PolygonAdapter) GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	// Polygon is a sidechain, not a rollup, so no L1 data fee
	return big.NewInt(0), nil
}

// IsDepositTransaction checks if a transaction is a deposit (always false for Polygon)
func (a *PolygonAdapter) IsDepositTransaction(tx *types.Transaction) bool {
	// Polygon doesn't have special deposit transaction types
	return false
}

// parseTransactionType is a helper to parse transaction type string
func parseTransactionType(typeStr string) TransactionType {
	if typeStr == "" {
		return TxTypeLegacy
	}
	
	txType, err := strconv.ParseInt(typeStr, 0, 64)
	if err != nil {
		return TxTypeLegacy
	}
	
	return TransactionType(txType)
}
