package etherscan

import (
	"context"
	"encoding/json"
	"net/url"
)

// Response with balance of a single address
type abiResponse struct {
	*baseResponse
	Data json.RawMessage `json:"result"`
}

// Parses a single balance response
func parseABIResponse(data []byte) ([]byte, error) {
	res := &abiResponse{baseResponse: &baseResponse{}}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// ContractABI returns the raw, unparsed ABI definition for a smart contract
// at the given address
func (c *Client) ContractABI(ctx context.Context, addr string) ([]byte, error) {
	params := url.Values{}
	params.Set("module", "contract")
	params.Set("action", "getabi")
	params.Set("address", addr)

	req, err := c.buildRequest(params)
	if err != nil {
		return nil, err
	}
	data, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseABIResponse(data)
}
