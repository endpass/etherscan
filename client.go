package etherscan

import (
	"fmt"
	"net/http"
	"strings"
)

// Client is the main client interface to the Etherscan API
type Client struct {
	// The root API endpoint, derived from the selected Network
	apiBase string

	// Ethereum network to get data from. Default: mainnet
	Network string

	// API Key to use for requests
	APIKey string

	// Wrapper *http.Client, can be replaced with your own client
	HTTPClient *http.Client
}

// Sets default values so that the zero value of *Client can be used
func (c *Client) setDefaults() error {
	if c.Network == "" {
		c.Network = "mainnet"
	}
	c.apiBase = apiEndpoints[strings.ToLower(c.Network)]
	// If still blank, invalid network
	if c.apiBase == "" {
		return fmt.Errorf("Invalid Network: %s. Network must be one of %s",
			c.Network, strings.Join(supportedNetworks(), ","))
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}
	return nil
}
