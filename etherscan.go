// Package etherscan is a client library for the etherscan.io API
package etherscan

import "time"

var (
	apiEndpoints = map[string]string{
		"mainnet": "https://api.etherscan.io/api",
		"ropsten": "https://api-ropsten.etherscan.io/api",
		"kovan":   "https://api-kovan.etherscan.io/api",
		"rinkeby": "https://api-rinkeby.etherscan.io/api",
	}

	userAgent     = "go-etherscan"
	clientTimeout = 30 * time.Second
)

// Returns supported networks based on API endpoints
func supportedNetworks() []string {
	var results []string
	for network := range apiEndpoints {
		results = append(results, network)
	}
	return results
}
