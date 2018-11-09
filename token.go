package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

// Token represents an ERC20 compatible token
type Token struct {
	Name   string
	Symbol string
	// Number of decimal places used by this token
	Decimals int
}

type tokenResponse struct {
	*baseResponse
	Total string `json:"result"`
}

func parseTokenResponse(r io.Reader) (*big.Int, error) {
	res := tokenResponse{baseResponse: &baseResponse{}}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}

	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}

	total := &big.Int{}
	total, ok := total.SetString(res.Total, 10)
	if !ok {
		return nil, errors.New("Could not parse total supply: " + res.Total)
	}

	return total, nil
}

func (c *Client) buildTokenTotalSupplyRequest(contractAddress string) (*http.Request, error) {
	if !strings.HasPrefix(contractAddress, "0x") {
		return nil, errors.New("Contract address must begin with 0x")
	}

	params := url.Values{}
	params.Set("module", "stats")
	params.Set("action", "tokensupply")
	params.Set("tag", "latest")
	params.Set("contractaddress", contractAddress)

	return c.buildRequest(params)
}

func (c *Client) buildTokenTotalBalanceRequest(contractAddress, address string) (*http.Request, error) {
	if !strings.HasPrefix(contractAddress, "0x") {
		return nil, errors.New("Contract address must begin with 0x")
	}

	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "tokenbalance")
	params.Set("tag", "latest")
	params.Set("contractaddress", contractAddress)
	params.Set("address", address)

	return c.buildRequest(params)
}

func (c *Client) tokenTotalSupply(ctx context.Context, contractAddress string) (*big.Int, error) {
	req, err := c.buildTokenTotalSupplyRequest(contractAddress)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseTokenResponse(resp.Body)
}

func (c *Client) tokenTotalBalance(ctx context.Context, contractAddress, address string) (*big.Int, error) {
	req, err := c.buildTokenTotalBalanceRequest(contractAddress, address)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	defer resp.Body.Close()
	return parseTokenResponse(resp.Body)
}

// TokenTotalSupply returns ERC20-Token TotalSupply by ContractAddress
func (c *Client) TokenTotalSupply(contractAddress string) (*big.Int, error) {
	return c.tokenTotalSupply(context.Background(), contractAddress)
}

// TokenTotalSupplyContext returns ERC20-Token TotalSupply by ContractAddress with a custom context
func (c *Client) TokenTotalSupplyContext(ctx context.Context, contractAddress string) (*big.Int, error) {
	return c.tokenTotalSupply(ctx, contractAddress)
}

// TokenTotalBalance returns ERC20-Token Account Balance for TokenContractAddress
func (c Client) TokenTotalBalance(contractAddress string, address string) (*big.Int, error) {
	return c.tokenTotalBalance(context.Background(), contractAddress, address)
}

// TokenTotalBalanceContext returns ERC20-Token Account Balance for TokenContractAddress with a custom context
func (c Client) TokenTotalBalanceContext(ctx context.Context, contractAddress string, address string) (*big.Int, error) {
	return c.tokenTotalBalance(ctx, contractAddress, address)
}
