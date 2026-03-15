package chain

import (
	"context"
	"math/big"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// TransactionType represents different transaction types across chains
type TransactionType int

const (
	TxTypeLegacy             TransactionType = 0   // Legacy transaction
	TxTypeAccessList         TransactionType = 1   // EIP-2930 with access list
	TxTypeDynamicFee         TransactionType = 2   // EIP-1559 with dynamic fee
	TxTypeArbitrumRetryable  TransactionType = 100 // Arbitrum retryable ticket
	TxTypeOptimismDeposit    TransactionType = 126 // Optimism deposit (0x7E)
)

// Adapter provides chain-specific operations for block processing and transaction handling
type Adapter interface {
	// Block operations
	NormalizeBlock(rawBlock *types.Block) (*types.Block, error)
	GetFinalityDepth() int64
	CalculateBlockReward(block *types.Block) (*big.Int, error)
	
	// Transaction operations
	NormalizeTransaction(tx *types.Transaction) (*types.Transaction, error)
	CalculateEffectiveFee(tx *types.Transaction, receipt *types.TransactionReceipt) (*big.Int, error)
	GetTransactionType(tx *types.Transaction) TransactionType
	
	// Chain-specific queries
	GetL1DataFee(ctx context.Context, tx *types.Transaction) (*big.Int, error)
	IsDepositTransaction(tx *types.Transaction) bool
}

// GetAdapter returns the appropriate adapter for a ChainID
func GetAdapter(chainID ChainID) Adapter {
	switch chainID {
	case EthereumMainnet, EthereumSepolia:
		return &EthereumAdapter{chainID: chainID}
	case ArbitrumMainnet, ArbitrumSepolia:
		return &ArbitrumAdapter{chainID: chainID}
	case OptimismMainnet, OptimismSepolia:
		return &OptimismAdapter{chainID: chainID}
	case AvalancheMainnet, AvalancheFuji:
		return &AvalancheAdapter{chainID: chainID}
	case PolygonMainnet, PolygonAmoy:
		return &PolygonAdapter{chainID: chainID}
	default:
		// Return Ethereum adapter as safe default
		return &EthereumAdapter{chainID: chainID}
	}
}
