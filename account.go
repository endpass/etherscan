package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/url"
)

// Response with balance of a single address
type balanceResponse struct {
	*baseResponse
	Balance string `json:"result"`
}

// Parses a single balance response
func parseBalanceResponse(data []byte) (*big.Int, error) {
	res := &balanceResponse{baseResponse: &baseResponse{}}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}

	bal := &big.Int{}
	var ok bool
	bal, ok = bal.SetString(res.Balance, 10)
	if !ok {
		return nil, errors.New("Could not parse balance: " + res.Balance)
	}
	return bal, nil
}

// Balance returns the balance of a single address
func (c *Client) Balance(ctx context.Context, addr string) (*big.Int, error) {
	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "balance")
	params.Set("tag", "latest")
	params.Set("address", "addr")

	req, err := c.buildRequest(params)
	if err != nil {
		return nil, err
	}
	data, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	return parseBalanceResponse(data)

}
