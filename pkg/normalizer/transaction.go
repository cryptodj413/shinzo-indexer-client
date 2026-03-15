package normalizer

import (
	"fmt"
	"strconv"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/chain"
	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// TransactionNormalizer handles chain-specific transaction normalization
type TransactionNormalizer struct {
	chainID chain.ChainID
	adapter chain.Adapter
}

// NewTransactionNormalizer creates a normalizer for a specific chain
func NewTransactionNormalizer(chainID chain.ChainID) *TransactionNormalizer {
	return &TransactionNormalizer{
		chainID: chainID,
		adapter: chain.GetAdapter(chainID),
	}
}

// Normalize converts raw transaction to normalized format
func (n *TransactionNormalizer) Normalize(rawTx *types.Transaction) (*types.Transaction, error) {
	if rawTx == nil {
		return nil, fmt.Errorf("raw transaction cannot be nil")
	}
	
	// Validate transaction structure
	if err := n.ValidateTransaction(rawTx); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}
	
	// Use chain adapter for normalization
	normalized, err := n.adapter.NormalizeTransaction(rawTx)
	if err != nil {
		return nil, fmt.Errorf("adapter normalization failed: %w", err)
	}
	
	// Validate transaction type is valid for this chain
	txType := n.adapter.GetTransactionType(normalized)
	if !n.isValidTypeForChain(txType) {
		return nil, fmt.Errorf("transaction type %d not valid for chain %d", txType, n.chainID)
	}
	
	// Validate type-specific fields
	if err := n.validateTypeSpecificFields(normalized, txType); err != nil {
		return nil, fmt.Errorf("type-specific validation failed: %w", err)
	}
	
	return normalized, nil
}

// ValidateTransaction ensures transaction meets basic requirements
func (n *TransactionNormalizer) ValidateTransaction(tx *types.Transaction) error {
	// Validate hash format (0x + 64 hex characters)
	if len(tx.Hash) != 66 {
		return fmt.Errorf("invalid transaction hash length: expected 66 characters (0x + 64 hex), got %d", len(tx.Hash))
	}
	if !hasHexPrefix(tx.Hash) {
		return fmt.Errorf("transaction hash must start with 0x")
	}
	
	return nil
}

// isValidTypeForChain checks if a transaction type is valid for the chain
func (n *TransactionNormalizer) isValidTypeForChain(txType chain.TransactionType) bool {
	switch n.chainID {
	case chain.EthereumMainnet, chain.EthereumSepolia:
		// Ethereum supports standard types only
		return txType == chain.TxTypeLegacy || txType == chain.TxTypeAccessList || txType == chain.TxTypeDynamicFee
	
	case chain.ArbitrumMainnet, chain.ArbitrumSepolia:
		// Arbitrum supports standard types plus retryable tickets
		return txType == chain.TxTypeLegacy || txType == chain.TxTypeAccessList || 
			txType == chain.TxTypeDynamicFee || txType == chain.TxTypeArbitrumRetryable
	
	case chain.OptimismMainnet, chain.OptimismSepolia:
		// Optimism supports standard types plus deposit transactions
		return txType == chain.TxTypeLegacy || txType == chain.TxTypeAccessList || 
			txType == chain.TxTypeDynamicFee || txType == chain.TxTypeOptimismDeposit
	
	case chain.AvalancheMainnet, chain.AvalancheFuji, chain.PolygonMainnet, chain.PolygonAmoy:
		// Standard EVM chains support standard types only
		return txType == chain.TxTypeLegacy || txType == chain.TxTypeAccessList || txType == chain.TxTypeDynamicFee
	
	default:
		// Unknown chain, allow standard types
		return txType == chain.TxTypeLegacy || txType == chain.TxTypeAccessList || txType == chain.TxTypeDynamicFee
	}
}

// validateTypeSpecificFields validates fields required for specific transaction types
func (n *TransactionNormalizer) validateTypeSpecificFields(tx *types.Transaction, txType chain.TransactionType) error {
	switch txType {
	case chain.TxTypeDynamicFee:
		// Type 2 requires EIP-1559 fields
		if tx.MaxFeePerGas == "" {
			return fmt.Errorf("Type 2 transaction missing maxFeePerGas")
		}
		if tx.MaxPriorityFeePerGas == "" {
			return fmt.Errorf("Type 2 transaction missing maxPriorityFeePerGas")
		}
		// Type 2 should have access list (can be empty)
		// AccessList field is optional, so no validation needed
	
	case chain.TxTypeAccessList:
		// Type 1 should have access list (can be empty)
		// AccessList field is optional, so no validation needed
	
	case chain.TxTypeArbitrumRetryable:
		// Arbitrum retryable tickets should have specific fields
		// These are optional in the struct, so no strict validation
		// The RPC layer will populate them if available
	
	case chain.TxTypeOptimismDeposit:
		// Optimism deposit transactions should have L1 fields
		// These are optional in the struct, so no strict validation
		// The RPC layer will populate them if available
	}
	
	return nil
}

// ValidateChainID validates that transaction chain ID matches expected chain
func ValidateChainID(tx *types.Transaction, expectedChainID chain.ChainID) error {
	if tx.ChainId == "" {
		// Pre-EIP-155 transactions don't have chain ID
		return nil
	}
	
	txChainID, err := strconv.ParseInt(tx.ChainId, 0, 64)
	if err != nil {
		return fmt.Errorf("invalid transaction chainId: %w", err)
	}
	
	if txChainID != int64(expectedChainID) {
		return fmt.Errorf("chain ID mismatch: expected %d, got %d", expectedChainID, txChainID)
	}
	
	return nil
}
