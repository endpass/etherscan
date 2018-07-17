package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Response with balance of a single address
type abiResponse struct {
	*baseResponse
	Data json.RawMessage `json:"result"`
}

// Parses a single balance response
func parseABIResponse(r io.Reader) ([]byte, error) {
	res := &abiResponse{baseResponse: &baseResponse{}}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}
	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) buildContractABIRequest(addr string) (*http.Request, error) {
	if !strings.HasPrefix(addr, "0x") {
		return nil, errors.New("Address must begin with 0x")
	}
	params := url.Values{}
	params.Set("module", "contract")
	params.Set("action", "getabi")
	params.Set("address", addr)

	return c.buildRequest(params)
}

func (c *Client) contractABI(ctx context.Context, addr string) ([]byte, error) {
	req, err := c.buildContractABIRequest(addr)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseABIResponse(resp.Body)
}

// ContractABI returns the raw, unparsed ABI definition for a smart contract
// at the given address
func (c *Client) ContractABI(addr string) ([]byte, error) {
	return c.contractABI(context.Background(), addr)
}

// ContractABI returns the raw, unparsed ABI definition for a smart contract
// at the given address with a custom context
func (c *Client) ContractABIContext(ctx context.Context, addr string) ([]byte, error) {
	return c.contractABI(ctx, addr)
}
