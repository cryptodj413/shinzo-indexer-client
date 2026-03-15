package chain

// ChainID is a type-safe identifier for supported blockchain networks
type ChainID int64

// Supported EVM chain IDs
const (
	EthereumMainnet  ChainID = 1
	EthereumSepolia  ChainID = 11155111
	ArbitrumMainnet  ChainID = 42161
	ArbitrumSepolia  ChainID = 421614
	OptimismMainnet  ChainID = 10
	OptimismSepolia  ChainID = 11155420
	AvalancheMainnet ChainID = 43114
	AvalancheFuji    ChainID = 43113
	PolygonMainnet   ChainID = 137
	PolygonAmoy      ChainID = 80002
)

// ChainInfo contains metadata about a blockchain network
type ChainInfo struct {
	ID          ChainID
	Name        string
	Network     string
	DisplayName string
	IsL2        bool
	ParentChain ChainID // For L2s, the L1 chain ID (0 for L1 chains)
}

// String returns the string representation of a ChainID
func (c ChainID) String() string {
	switch c {
	case EthereumMainnet:
		return "Ethereum Mainnet"
	case EthereumSepolia:
		return "Ethereum Sepolia"
	case ArbitrumMainnet:
		return "Arbitrum Mainnet"
	case ArbitrumSepolia:
		return "Arbitrum Sepolia"
	case OptimismMainnet:
		return "Optimism Mainnet"
	case OptimismSepolia:
		return "Optimism Sepolia"
	case AvalancheMainnet:
		return "Avalanche Mainnet"
	case AvalancheFuji:
		return "Avalanche Fuji"
	case PolygonMainnet:
		return "Polygon Mainnet"
	case PolygonAmoy:
		return "Polygon Amoy"
	default:
		return "Unknown Chain"
	}
}
