package etherscan

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
		c.HTTPClient = &http.Client{
			Timeout: clientTimeout,
		}
	}
	return nil
}

// Construct a new request to the API that is ready to send
// All methods use GET requests for now
func (c *Client) buildRequest(params url.Values) (*http.Request, error) {
	if err := c.setDefaults(); err != nil {
		return nil, err
	}

	if params == nil {
		return nil, errors.New("Params are empty")
	}
	if params.Get("module") == "" {
		return nil, errors.New("Missing required parameter: module")
	}
	if params.Get("action") == "" {
		return nil, errors.New("Missing required parameter: action")
	}
	if params.Get("apikey") == "" {
		params.Set("apikey", c.APIKey)
	}

	reqURL := c.apiBase + params.Encode()
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

// Sends a request and returns the response body
func (c *Client) sendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("Context is nil")
	}
	if req == nil {
		return nil, errors.New("Request is nil")
	}
	req = req.WithContext(ctx)
	return c.HTTPClient.Do(req)
}
