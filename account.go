package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/url"
	"strings"
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

func (c *Client) buildBalanceRequest(addr string) (*http.Request, error) {
	if !strings.HasPrefix(addr, "0x") {
		return nil, errors.New("Address must begin with 0x")
	}
	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "balance")
	params.Set("tag", "latest")
	params.Set("address", addr)

	return c.buildRequest(params)
}

func (c *Client) balance(ctx context.Context, addr string) (*big.Int, error) {
	req, err := c.buildBalanceRequest(addr)
	if err != nil {
		return nil, err
	}
	data, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return parseBalanceResponse(data)
}

// Balance returns the balance of a single address
func (c *Client) Balance(addr string) (*big.Int, error) {
	return c.balance(context.Background(), addr)
}

// BalanceContext returns the balance of a single address with a custom
// context
func (c *Client) BalanceContext(ctx context.Context, addr string) (*big.Int, error) {
	return c.balance(ctx, addr)
}
