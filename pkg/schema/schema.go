package schema

import (
	_ "embed"
	"strings"
	
	"github.com/shinzonetwork/shinzo-indexer-client/pkg/chain"
)

//go:embed schema_standard.graphql
var SchemaGraphQL string

//go:embed schema_branchable.graphql
var SchemaBranchGraphQL string

// GetSchema returns the GraphQL schema found in `schema.graphql` as a string.
func GetSchema() string {
	return SchemaGraphQL
}

// GetBranchableSchema returns the branchable GraphQL schema.
func GetBranchableSchema() string {
	return SchemaBranchGraphQL
}

// GetSchemaForBuild returns the appropriate schema based on build tags.
func GetSchemaForBuild() string {
	return schemaForBuild(IsBranchable())
}

func schemaForBuild(branchable bool) string {
	if branchable {
		return GetBranchableSchema()
	}
	return GetSchema()
}

// GetSchemaForChain returns the schema with collection names adapted for the given chain prefix.
// It uses the schema manager to generate chain-specific schemas (L1 vs L2).
// TODO: This function currently uses a simple prefix replacement approach for backward compatibility.
// Future enhancement: Use Manager.GetSchemaForChain() to generate proper L1/L2 schemas based on ChainID.
func GetSchemaForChain(prefix string) string {
	s := GetSchemaForBuild()
	if prefix == "" || prefix == "Ethereum__Mainnet__" {
		return s
	}
	// Remove trailing __ from prefix if present for replacement
	cleanPrefix := strings.TrimSuffix(prefix, "__")
	return strings.ReplaceAll(s, "Ethereum__Mainnet", cleanPrefix)
}

// GetSchemaForChainID returns the schema for a specific chain using the Manager.
// This generates proper L1/L2-specific schemas based on the chain type.
// TODO: Replace GetSchemaForChain() with this function once all callers are updated.
func GetSchemaForChainID(chainID chain.ChainID) (string, error) {
	registry := chain.NewRegistry()
	manager := NewManager(registry)
	return manager.GetSchemaForChain(chainID)
}
