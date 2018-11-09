package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
)

type statsTotalSupplyResponse struct {
	*baseResponse
	Total string `json:"result"`
}

type statsLastPriceResposne struct {
	*baseResponse
	LastPrice *lastPrice `json:"result"`
}

type lastPrice struct {
	Ethbtc          string `json:"ethbtc"`
	EthbtcTimestamp string `json:"ethbtc_timestamp"`
	Ethusd          string `json:"ethusd"`
	EthusdTimestamp string `json:"ethusd_timestamp"`
}

type LastPrice struct {
	Ethbtc          *big.Float
	EthbtcTimestamp int
	Ethusd          *big.Float
	EthusdTimestamp int
}

func parseStatsTotalSupplyResponse(r io.Reader) (*big.Int, error) {
	res := statsTotalSupplyResponse{baseResponse: &baseResponse{}}
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

func parseStatsLastPriceResponse(r io.Reader) (*LastPrice, error) {
	res := statsLastPriceResposne{baseResponse: &baseResponse{}}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}

	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}

	if res.LastPrice == nil {
		return nil, errors.New("result is empty")
	}

	lp := &LastPrice{
		Ethbtc:          parseFloat(res.LastPrice.Ethbtc),
		EthbtcTimestamp: parseInt(res.LastPrice.EthbtcTimestamp),
		Ethusd:          parseFloat(res.LastPrice.Ethusd),
		EthusdTimestamp: parseInt(res.LastPrice.EthusdTimestamp),
	}

	return lp, nil
}

func (c *Client) buildStatsTotalSupplyRequest() (*http.Request, error) {
	params := url.Values{}
	params.Set("module", "stats")
	params.Set("action", "ethsupply")

	return c.buildRequest(params)
}

func (c *Client) buildStatsLastPriceRequest() (*http.Request, error) {
	params := url.Values{}
	params.Set("module", "stats")
	params.Set("action", "ethprice")
	return c.buildRequest(params)
}

func (c *Client) statsTotalSupply(ctx context.Context) (*big.Int, error) {
	req, err := c.buildStatsTotalSupplyRequest()
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseStatsTotalSupplyResponse(resp.Body)
}

func (c *Client) statsLastPrice(ctx context.Context) (*LastPrice, error) {
	req, err := c.buildStatsLastPriceRequest()
	if err != nil {
		return nil, err
	}
	fmt.Println(req.URL.String())
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseStatsLastPriceResponse(resp.Body)
}

// TotalSupply returns total supply of ether
func (c *Client) TotalSupply() (*big.Int, error) {
	return c.statsTotalSupply(context.Background())
}

// TotalSupplyContext returns total supply of ether with a custom context
func (c *Client) TotalSupplyContext(ctx context.Context) (*big.Int, error) {
	return c.statsTotalSupply(ctx)
}

// LastPrice returns ETHER last price
func (c *Client) LastPrice() (*LastPrice, error) {
	return c.statsLastPrice(context.Background())
}

// LastPriceContext returns ETHER last price with a custom context
func (c *Client) LastPriceContext(ctx context.Context) (*LastPrice, error) {
	return c.statsLastPrice(ctx)
}
