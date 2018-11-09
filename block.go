package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
)

// Block is a single block in the chain
type Block struct {
	Number int
	Hash   string
	// All transactions mined in this block
	Transactions []*Transaction
}

// Response with block and uncle rewards
type blockResponse struct {
	*baseResponse
	BlockReward *blockRewardResponse `json:"result"`
}

// Unparsed block
type blockRewardResponse struct {
	BlockNumber          string               `json:"blockNumber"`
	TimeStamp            string               `json:"timeStamp"`
	BlockMiner           string               `json:"blockMiner"`
	BlockReward          string               `json:"blockReward"`
	Uncles               []blockUncleResponse `json:"uncles"`
	UncleInclusionReward string               `json:"uncleInclusionReward"`
}

// Unparsed block uncle
type blockUncleResponse struct {
	Miner         string `json:"miner"`
	UnclePosition string `json:"unclePosition"`
	Blockreward   string `json:"blockreward"`
}

type BlockReward struct {
	BlockNumber          int
	TimeStamp            int
	BlockMiner           string
	BlockReward          *big.Int
	Uncles               []BlockUncle
	UncleInclusionReward *big.Int
}

type BlockUncle struct {
	Miner         string
	UnclePosition int
	BlockReward   *big.Int
}

func parseBlockRewardResponse(r io.Reader) (*BlockReward, error) {
	res := blockResponse{baseResponse: &baseResponse{}}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}

	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}

	if res.BlockReward == nil {
		return nil, errors.New("result is empty")
	}

	blockInfo, err := parseBlock(res.BlockReward)
	if err != nil {
		return nil, err
	}

	return blockInfo, nil
}

func parseBlock(blockResponse *blockRewardResponse) (*BlockReward, error) {
	block := &BlockReward{
		BlockNumber:          parseInt(blockResponse.BlockNumber),
		TimeStamp:            parseInt(blockResponse.TimeStamp),
		BlockMiner:           blockResponse.BlockMiner,
		BlockReward:          parseBig(blockResponse.BlockReward),
		UncleInclusionReward: parseBig(blockResponse.UncleInclusionReward),
	}

	uncles := make([]BlockUncle, len(blockResponse.Uncles))
	for i, u := range blockResponse.Uncles {
		uncles[i] = BlockUncle{
			Miner:         u.Miner,
			UnclePosition: parseInt(u.UnclePosition),
			BlockReward:   parseBig(u.Blockreward),
		}
	}
	block.Uncles = uncles

	return block, nil
}

func (c *Client) buildBlockRequest(blockNumber int) (*http.Request, error) {
	params := url.Values{}
	params.Set("module", "block")
	params.Set("action", "getblockreward")
	params.Set("blockno", strconv.Itoa(blockNumber))

	return c.buildRequest(params)
}

func (c *Client) blockReward(ctx context.Context, blockNumber int) (*BlockReward, error) {
	req, err := c.buildBlockRequest(blockNumber)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseBlockRewardResponse(resp.Body)
}

// BlockReward returns the reward of single block
func (c *Client) BlockReward(blockNumber int) (*BlockReward, error) {
	return c.blockReward(context.Background(), blockNumber)
}

// BlockRewardContext returns the reward of single block with a custom context
func (c *Client) BlockRewardContext(ctx context.Context, blockNumber int) (*BlockReward, error) {
	return c.blockReward(ctx, blockNumber)
}
