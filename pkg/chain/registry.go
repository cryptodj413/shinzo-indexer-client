package chain

import (
	"fmt"
	"strings"
)

// Registry manages chain validation and lookup
type Registry struct {
	chains map[string]ChainInfo // key: "name__network" (lowercase)
	byID   map[ChainID]ChainInfo
}

// NewRegistry creates a new chain registry with all supported chains
func NewRegistry() *Registry {
	r := &Registry{
		chains: make(map[string]ChainInfo),
		byID:   make(map[ChainID]ChainInfo),
	}
	r.registerChains()
	return r
}

// registerChains populates the registry with all supported chains
func (r *Registry) registerChains() {
	chains := []ChainInfo{
		{
			ID:          EthereumMainnet,
			Name:        "Ethereum",
			Network:     "Mainnet",
			DisplayName: "Ethereum Mainnet",
			IsL2:        false,
			ParentChain: 0,
		},
		{
			ID:          EthereumSepolia,
			Name:        "Ethereum",
			Network:     "Sepolia",
			DisplayName: "Ethereum Sepolia Testnet",
			IsL2:        false,
			ParentChain: 0,
		},
		{
			ID:          ArbitrumMainnet,
			Name:        "Arbitrum",
			Network:     "Mainnet",
			DisplayName: "Arbitrum One",
			IsL2:        true,
			ParentChain: EthereumMainnet,
		},
		{
			ID:          ArbitrumSepolia,
			Name:        "Arbitrum",
			Network:     "Sepolia",
			DisplayName: "Arbitrum Sepolia Testnet",
			IsL2:        true,
			ParentChain: EthereumSepolia,
		},
		{
			ID:          OptimismMainnet,
			Name:        "Optimism",
			Network:     "Mainnet",
			DisplayName: "Optimism Mainnet",
			IsL2:        true,
			ParentChain: EthereumMainnet,
		},
		{
			ID:          OptimismSepolia,
			Name:        "Optimism",
			Network:     "Sepolia",
			DisplayName: "Optimism Sepolia Testnet",
			IsL2:        true,
			ParentChain: EthereumSepolia,
		},
		{
			ID:          AvalancheMainnet,
			Name:        "Avalanche",
			Network:     "Mainnet",
			DisplayName: "Avalanche C-Chain",
			IsL2:        false,
			ParentChain: 0,
		},
		{
			ID:          AvalancheFuji,
			Name:        "Avalanche",
			Network:     "Fuji",
			DisplayName: "Avalanche Fuji Testnet",
			IsL2:        false,
			ParentChain: 0,
		},
		{
			ID:          PolygonMainnet,
			Name:        "Polygon",
			Network:     "Mainnet",
			DisplayName: "Polygon PoS",
			IsL2:        false,
			ParentChain: 0,
		},
		{
			ID:          PolygonAmoy,
			Name:        "Polygon",
			Network:     "Amoy",
			DisplayName: "Polygon Amoy Testnet",
			IsL2:        false,
			ParentChain: 0,
		},
	}

	for _, info := range chains {
		key := r.makeKey(info.Name, info.Network)
		r.chains[key] = info
		r.byID[info.ID] = info
	}
}

// makeKey creates a normalized lookup key from name and network
func (r *Registry) makeKey(name, network string) string {
	return strings.ToLower(name) + "__" + strings.ToLower(network)
}

// ValidateAndGetChainID validates user input and returns the corresponding ChainID
func (r *Registry) ValidateAndGetChainID(name, network string) (ChainID, error) {
	if strings.TrimSpace(name) == "" {
		return 0, fmt.Errorf("chain name cannot be empty")
	}
	if strings.TrimSpace(network) == "" {
		return 0, fmt.Errorf("chain network cannot be empty")
	}

	key := r.makeKey(name, network)
	info, exists := r.chains[key]
	if !exists {
		return 0, r.unsupportedChainError(name, network)
	}

	return info.ID, nil
}

// GetChainInfo returns metadata for a ChainID
func (r *Registry) GetChainInfo(id ChainID) (ChainInfo, error) {
	info, exists := r.byID[id]
	if !exists {
		return ChainInfo{}, fmt.Errorf("chain ID %d not found in registry", id)
	}
	return info, nil
}

// IsSupported checks if a chain name and network combination is supported
func (r *Registry) IsSupported(name, network string) bool {
	key := r.makeKey(name, network)
	_, exists := r.chains[key]
	return exists
}

// ListSupportedChains returns all supported chains
func (r *Registry) ListSupportedChains() []ChainInfo {
	chains := make([]ChainInfo, 0, len(r.byID))
	for _, info := range r.byID {
		chains = append(chains, info)
	}
	return chains
}

// unsupportedChainError creates a detailed error message for unsupported chains
func (r *Registry) unsupportedChainError(name, network string) error {
	supported := make(map[string][]string)
	for _, info := range r.byID {
		supported[info.Name] = append(supported[info.Name], info.Network)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("unsupported chain: %s %s\n", name, network))
	sb.WriteString("Supported chains:\n")
	for chainName, networks := range supported {
		sb.WriteString(fmt.Sprintf("  - %s: %s\n", chainName, strings.Join(networks, ", ")))
	}

	return fmt.Errorf("%s", sb.String())
}
