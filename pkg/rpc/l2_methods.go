package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/logger"
	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// GetL1DataFee attempts to get the L1 data fee for an L2 transaction
// Falls back to estimation if the RPC method is not available
func (c *EthereumClient) GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error) {
	// If already populated in transaction, return it
	if tx.L1DataFee != "" {
		l1DataFee := new(big.Int)
		if _, ok := l1DataFee.SetString(tx.L1DataFee, 0); !ok {
			return nil, fmt.Errorf("invalid l1DataFee: %s", tx.L1DataFee)
		}
		return l1DataFee, nil
	}
	
	// Try to call chain-specific RPC method
	// This would require chain-specific RPC calls which may not be available
	// For now, estimate based on transaction data size
	l1DataFee := estimateL1DataFee(tx)
	
	logger.Sugar.Debugf("L1 data fee method not available for tx %s, using estimate: %s", 
		tx.Hash, l1DataFee.String())
	
	return l1DataFee, nil
}

// estimateL1DataFee estimates the L1 data fee based on transaction data size
// This is a rough approximation when the actual RPC method is not available
func estimateL1DataFee(tx *types.Transaction) *big.Int {
	// Rough estimation based on input data size
	// L1 data fee is typically proportional to the calldata size
	// This is a placeholder - actual calculation would need L1 gas price and compression ratio
	
	inputSize := len(tx.Input)
	if inputSize < 2 {
		// No meaningful data
		return big.NewInt(0)
	}
	
	// Remove "0x" prefix for size calculation
	dataSize := (inputSize - 2) / 2 // Convert hex string to byte count
	
	// Rough estimate: ~16 gas per byte of calldata on L1
	// Multiply by an estimated L1 gas price (this is very rough)
	estimatedL1Gas := int64(dataSize * 16)
	estimatedL1GasPrice := int64(20_000_000_000) // 20 gwei estimate
	
	estimate := big.NewInt(estimatedL1Gas)
	estimate.Mul(estimate, big.NewInt(estimatedL1GasPrice))
	
	return estimate
}

// GetArbitrumL1DataFee attempts to get L1 data fee using Arbitrum-specific RPC method
func (c *EthereumClient) GetArbitrumL1DataFee(ctx context.Context, txHash string) (*big.Int, error) {
	// Arbitrum-specific RPC method would be called here
	// For now, return 0 and let the estimation happen in GetL1DataFee
	logger.Sugar.Debugf("Arbitrum L1 data fee RPC method not implemented, using estimation")
	return big.NewInt(0), nil
}

// GetOptimismL1DataFee attempts to get L1 data fee using Optimism-specific RPC method
func (c *EthereumClient) GetOptimismL1DataFee(ctx context.Context, txHash string) (*big.Int, error) {
	// Optimism-specific RPC method would be called here
	// For now, return 0 and let the estimation happen in GetL1DataFee
	logger.Sugar.Debugf("Optimism L1 data fee RPC method not implemented, using estimation")
	return big.NewInt(0), nil
}

// IsMethodAvailable checks if a specific RPC method is available on the node
func (c *EthereumClient) IsMethodAvailable(ctx context.Context, method string) bool {
	// This would test if the RPC method exists
	// For now, assume standard methods are available
	standardMethods := map[string]bool{
		"eth_blockNumber":        true,
		"eth_getBlockByNumber":   true,
		"eth_getTransactionReceipt": true,
		"eth_getBlockReceipts":   true,
	}
	
	return standardMethods[method]
}
