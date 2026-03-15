package normalizer

import (
	"fmt"
	"strconv"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/chain"
	"github.com/shinzonetwork/shinzo-indexer-client/pkg/types"
)

// BlockNormalizer handles chain-specific block normalization
type BlockNormalizer struct {
	chainID chain.ChainID
	adapter chain.Adapter
}

// NewBlockNormalizer creates a normalizer for a specific chain
func NewBlockNormalizer(chainID chain.ChainID) *BlockNormalizer {
	return &BlockNormalizer{
		chainID: chainID,
		adapter: chain.GetAdapter(chainID),
	}
}

// Normalize converts raw block to normalized format
func (n *BlockNormalizer) Normalize(rawBlock *types.Block) (*types.Block, error) {
	if rawBlock == nil {
		return nil, fmt.Errorf("raw block cannot be nil")
	}
	
	// Validate block structure
	if err := n.ValidateBlock(rawBlock); err != nil {
		return nil, fmt.Errorf("block validation failed: %w", err)
	}
	
	// Use chain adapter for normalization
	normalized, err := n.adapter.NormalizeBlock(rawBlock)
	if err != nil {
		return nil, fmt.Errorf("adapter normalization failed: %w", err)
	}
	
	// Set optional fields based on chain
	if err := n.SetOptionalFields(normalized); err != nil {
		return nil, fmt.Errorf("failed to set optional fields: %w", err)
	}
	
	return normalized, nil
}

// SetOptionalFields sets chain-specific optional fields
func (n *BlockNormalizer) SetOptionalFields(block *types.Block) error {
	// Chain adapter already handles this in NormalizeBlock
	// This method is for any additional post-processing if needed
	return nil
}

// ValidateBlock ensures block meets chain-specific requirements
func (n *BlockNormalizer) ValidateBlock(block *types.Block) error {
	// Validate hash format (0x + 64 hex characters)
	if len(block.Hash) != 66 {
		return fmt.Errorf("invalid block hash length: expected 66 characters (0x + 64 hex), got %d", len(block.Hash))
	}
	if !hasHexPrefix(block.Hash) {
		return fmt.Errorf("block hash must start with 0x")
	}
	
	// Validate block number is non-negative
	blockNum, err := strconv.ParseInt(block.Number, 0, 64)
	if err != nil {
		return fmt.Errorf("invalid block number: %w", err)
	}
	if blockNum < 0 {
		return fmt.Errorf("block number must be non-negative, got %d", blockNum)
	}
	
	// Validate timestamp
	timestamp, err := strconv.ParseInt(block.Timestamp, 0, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %w", err)
	}
	if timestamp < 0 {
		return fmt.Errorf("timestamp must be non-negative, got %d", timestamp)
	}
	
	return nil
}
