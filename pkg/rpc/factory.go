package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/chain"
)

// ClientConfig holds RPC configuration for a chain
type ClientConfig struct {
	NodeURL string
	WsURL   string
	APIKey  string
	Timeout time.Duration
}

// ClientFactory creates RPC clients for different chains
type ClientFactory struct {
	configs map[chain.ChainID]ClientConfig
}

// NewClientFactory creates a factory with configuration
func NewClientFactory(configs map[chain.ChainID]ClientConfig) *ClientFactory {
	return &ClientFactory{
		configs: configs,
	}
}

// GetClient returns an RPC client for the specified chain
func (f *ClientFactory) GetClient(chainID chain.ChainID) (*EthereumClient, error) {
	config, exists := f.configs[chainID]
	if !exists {
		return nil, fmt.Errorf("no RPC configuration found for chain ID %d", chainID)
	}
	
	// Create client with configuration
	client, err := NewEthereumClient(config.NodeURL, config.WsURL, config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client for chain %d: %w", chainID, err)
	}
	
	return client, nil
}

// ValidateConnection tests the RPC connection for a chain
func (f *ClientFactory) ValidateConnection(chainID chain.ChainID) error {
	client, err := f.GetClient(chainID)
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	defer client.Close()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Try to get the latest block number as a connectivity test
	_, err = client.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("RPC connection validation failed: %w", err)
	}
	
	return nil
}

// AddConfig adds or updates RPC configuration for a chain
func (f *ClientFactory) AddConfig(chainID chain.ChainID, config ClientConfig) {
	f.configs[chainID] = config
}

// GetConfig returns the RPC configuration for a chain
func (f *ClientFactory) GetConfig(chainID chain.ChainID) (ClientConfig, bool) {
	config, exists := f.configs[chainID]
	return config, exists
}
