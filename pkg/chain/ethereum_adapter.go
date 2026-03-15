package chain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// EthereumAdapter implements the Adapter interface for Ethereum chains
type EthereumAdapter struct {
	chainID ChainID
}

// NormalizeBlock normalizes an Ethereum block, including PoW fields for pre-merge blocks
func (a *EthereumAdapter) NormalizeBlock(rawBlock *types.Block) (*types.Block, error) {
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
	
	// Ethereum blocks include all PoW fields
	// For post-merge blocks, these fields will be present but may have default values
	normalized := *rawBlock
	return &normalized, nil
}

// GetFinalityDepth returns the number of blocks for finality on Ethereum
func (a *EthereumAdapter) GetFinalityDepth() int64 {
	// Ethereum uses 32 blocks for finality (2 epochs in PoS)
	return 32
}

// CalculateBlockReward calculates the block reward for Ethereum
func (a *EthereumAdapter) CalculateBlockReward(block *types.Block) (*big.Int, error) {
	// Ethereum block rewards vary by era
	// Post-merge: No block reward (PoS validators get fees)
	// This is a simplified implementation
	return big.NewInt(0), nil
}

// NormalizeTransaction normalizes an Ethereum transaction
func (a *EthereumAdapter) NormalizeTransaction(tx *types.Transaction) (*types.Transaction, error) {
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

// CalculateEffectiveFee calculates the total fee for an Ethereum transaction
func (a *EthereumAdapter) CalculateEffectiveFee(tx *types.Transaction, receipt *types.TransactionReceipt) (*big.Int, error) {
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
		// Fallback to gasPrice for legacy transactions
		if _, ok := effectiveGasPrice.SetString(tx.GasPrice, 0); !ok {
			return nil, fmt.Errorf("invalid gasPrice: %s", tx.GasPrice)
		}
	}
	
	// Total fee = effectiveGasPrice * gasUsed
	totalFee := new(big.Int).Mul(effectiveGasPrice, gasUsed)
	return totalFee, nil
}

// GetTransactionType identifies the transaction type
func (a *EthereumAdapter) GetTransactionType(tx *types.Transaction) TransactionType {
	if tx.Type == "" {
		return TxTypeLegacy
	}
	
	txType, err := strconv.ParseInt(tx.Type, 0, 64)
	if err != nil {
		return TxTypeLegacy
	}
	
	return TransactionType(txType)
}

// GetL1DataFee returns the L1 data fee (always 0 for Ethereum L1)
func (a *EthereumAdapter) GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	// Ethereum is L1, no L1 data fee
	return big.NewInt(0), nil
}

// IsDepositTransaction checks if a transaction is a deposit (always false for Ethereum L1)
func (a *EthereumAdapter) IsDepositTransaction(tx *types.Transaction) bool {
	// Ethereum L1 has no deposit transactions
	return false
}

// Validation helpers

func validateBlockHash(hash string) error {
	if len(hash) != 66 {
		return fmt.Errorf("invalid block hash length: expected 66 characters (0x + 64 hex), got %d", len(hash))
	}
	if !hasHexPrefix(hash) {
		return fmt.Errorf("block hash must start with 0x")
	}
	return nil
}

func validateBlockNumber(number string) error {
	num, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return fmt.Errorf("invalid block number: %w", err)
	}
	if num < 0 {
		return fmt.Errorf("block number must be non-negative, got %d", num)
	}
	return nil
}

func validateTransactionHash(hash string) error {
	if len(hash) != 66 {
		return fmt.Errorf("invalid transaction hash length: expected 66 characters (0x + 64 hex), got %d", len(hash))
	}
	if !hasHexPrefix(hash) {
		return fmt.Errorf("transaction hash must start with 0x")
	}
	return nil
}

func hasHexPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}
