package schema

import (
	"fmt"
	"strings"

	"github.com/shinzonetwork/shinzo-indexer-client/pkg/chain"
)

// Manager handles schema generation and application for different chains
type Manager struct {
	registry *chain.Registry
}

// NewManager creates a schema manager
func NewManager(registry *chain.Registry) *Manager {
	return &Manager{
		registry: registry,
	}
}

// GetCollectionPrefix returns the collection name prefix for a chain
func (m *Manager) GetCollectionPrefix(chainID chain.ChainID) (string, error) {
	info, err := m.registry.GetChainInfo(chainID)
	if err != nil {
		return "", fmt.Errorf("failed to get chain info: %w", err)
	}
	
	// Format: {ChainName}__{Network}__
	prefix := fmt.Sprintf("%s__%s__", info.Name, info.Network)
	return prefix, nil
}

// GetSchemaForChain returns GraphQL schema for a specific chain
func (m *Manager) GetSchemaForChain(chainID chain.ChainID) (string, error) {
	prefix, err := m.GetCollectionPrefix(chainID)
	if err != nil {
		return "", err
	}
	
	info, err := m.registry.GetChainInfo(chainID)
	if err != nil {
		return "", err
	}
	
	// Generate schema based on chain type
	if info.IsL2 {
		return m.generateL2Schema(prefix, info)
	}
	return m.generateL1Schema(prefix, info)
}

// generateL1Schema generates schema for L1 chains (Ethereum, Avalanche, Polygon)
func (m *Manager) generateL1Schema(prefix string, info chain.ChainInfo) (string, error) {
	schema := fmt.Sprintf(`type %sBlock {
    hash: String
    number: Int
    timestamp: String
    parentHash: String
    difficulty: String
    totalDifficulty: String
    gasUsed: String
    gasLimit: String
    baseFeePerGas: String
    nonce: String
    miner: String
    size: String
    stateRoot: String
    sha3Uncles: String
    transactionsRoot: String
    receiptsRoot: String
    logsBloom: String
    extraData: String
    mixHash: String
    uncles: [String]
    transactions: [%sTransaction] @relation(name: "block_transactions")
}

type %sBlockSignature {
    blockNumber: Int
    blockHash: String
    merkleRoot: String
    cidCount: Int
    cids: [String]
    signatureType: String
    signatureIdentity: String
    signatureValue: String
    createdAt: String
}

type %sSnapshotSignature {
    startBlock: Int
    endBlock: Int
    merkleRoot: String
    blockCount: Int
    signatureType: String
    signatureIdentity: String
    signatureValue: String
    snapshotFile: String
    createdAt: String
    blockSigMerkleRoots: [String]
}

type %sTransaction {
    hash: String
    blockHash: String
    blockNumber: Int
    from: String
    to: String
    value: String
    gas: String
    gasPrice: String
    gasUsed: String
    maxFeePerGas: String
    maxPriorityFeePerGas: String
    input: String
    nonce: String
    transactionIndex: Int
    type: String
    chainId: String
    v: String
    r: String
    s: String
    status: Boolean
    cumulativeGasUsed: String
    effectiveGasPrice: String
    block: %sBlock @relation(name: "block_transactions")
    logs: [%sLog] @relation(name: "transaction_logs")
    accessList: [%sAccessListEntry] @relation(name: "transaction_accessList")
}

type %sAccessListEntry {
    address: String
    blockNumber: Int
    storageKeys: [String]
    transaction: %sTransaction @relation(name: "transaction_accessList")
}

type %sLog {
    address: String
    topics: [String]
    data: String
    transactionHash: String
    blockHash: String
    blockNumber: Int
    transactionIndex: Int
    logIndex: Int
    removed: String
    block: %sBlock @relation(name: "block_logs")
    transaction: %sTransaction @relation(name: "transaction_logs")
}
`, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix)
	
	return schema, nil
}

// generateL2Schema generates schema for L2 chains (Arbitrum, Optimism)
func (m *Manager) generateL2Schema(prefix string, info chain.ChainInfo) (string, error) {
	schema := fmt.Sprintf(`type %sBlock {
    hash: String
    number: Int
    timestamp: String
    parentHash: String
    gasUsed: String
    gasLimit: String
    baseFeePerGas: String
    miner: String
    size: String
    stateRoot: String
    transactionsRoot: String
    receiptsRoot: String
    logsBloom: String
    extraData: String
    l1BlockNumber: String
    sendCount: String
    sendRoot: String
    transactions: [%sTransaction] @relation(name: "block_transactions")
}

type %sBlockSignature {
    blockNumber: Int
    blockHash: String
    merkleRoot: String
    cidCount: Int
    cids: [String]
    signatureType: String
    signatureIdentity: String
    signatureValue: String
    createdAt: String
}

type %sSnapshotSignature {
    startBlock: Int
    endBlock: Int
    merkleRoot: String
    blockCount: Int
    signatureType: String
    signatureIdentity: String
    signatureValue: String
    snapshotFile: String
    createdAt: String
    blockSigMerkleRoots: [String]
}

type %sTransaction {
    hash: String
    blockHash: String
    blockNumber: Int
    from: String
    to: String
    value: String
    gas: String
    gasPrice: String
    gasUsed: String
    maxFeePerGas: String
    maxPriorityFeePerGas: String
    input: String
    nonce: String
    transactionIndex: Int
    type: String
    chainId: String
    v: String
    r: String
    s: String
    status: Boolean
    cumulativeGasUsed: String
    effectiveGasPrice: String
    l1BlockNumber: String
    l1Timestamp: String
    l1DataFee: String
    queueOrigin: String
    retryTo: String
    retryValue: String
    retryData: String
    beneficiary: String
    maxSubmissionFee: String
    block: %sBlock @relation(name: "block_transactions")
    logs: [%sLog] @relation(name: "transaction_logs")
    accessList: [%sAccessListEntry] @relation(name: "transaction_accessList")
}

type %sAccessListEntry {
    address: String
    blockNumber: Int
    storageKeys: [String]
    transaction: %sTransaction @relation(name: "transaction_accessList")
}

type %sLog {
    address: String
    topics: [String]
    data: String
    transactionHash: String
    blockHash: String
    blockNumber: Int
    transactionIndex: Int
    logIndex: Int
    removed: String
    block: %sBlock @relation(name: "block_logs")
    transaction: %sTransaction @relation(name: "transaction_logs")
}
`, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix)
	
	return schema, nil
}

// ApplySchema applies schema to DefraDB (placeholder for actual implementation)
func (m *Manager) ApplySchema(chainID chain.ChainID, defraURL string) error {
	schema, err := m.GetSchemaForChain(chainID)
	if err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}
	
	// This would call the actual DefraDB schema application logic
	// For now, just validate the schema is not empty
	if strings.TrimSpace(schema) == "" {
		return fmt.Errorf("generated schema is empty")
	}
	
	return nil
}

// ValidateCollectionName validates that a collection name matches the expected chain prefix
func ValidateCollectionName(collectionName string, chainID chain.ChainID, registry *chain.Registry) error {
	manager := NewManager(registry)
	expectedPrefix, err := manager.GetCollectionPrefix(chainID)
	if err != nil {
		return fmt.Errorf("failed to get expected prefix: %w", err)
	}
	
	if !strings.HasPrefix(collectionName, expectedPrefix) {
		return fmt.Errorf("collection name %s does not match chain %d prefix %s",
			collectionName, chainID, expectedPrefix)
	}
	
	return nil
}
