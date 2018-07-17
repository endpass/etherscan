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
	"time"
)

// Response with list of transactions for an address
type transactionsResponse struct {
	*baseResponse
	Transactions []*transactionResponse `json:"result"`
}

// An unparsed transaction
type transactionResponse struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

// Transaction represents a standard Ethereum transaction
type Transaction struct {
	Block *Block

	Timestamp time.Time

	Hash string

	Nonce int

	// Index of this transaction in the block
	Index int

	// From address hex
	From string

	To string

	// Address of contract for this transaction, if any
	ContractAddress string

	// Amount of transaction in wei
	Value *big.Int

	GasLimit int

	GasUsed int

	// Gas price in wei
	GasPrice *big.Int

	IsError bool

	Confirmations uint64

	// Data sent to transaction, encoded as hex
	Data string
}

func parseTransaction(tx *transactionResponse) *Transaction {
	return &Transaction{
		Block: &Block{
			Number: parseInt(tx.BlockNumber),
			Hash:   tx.BlockHash,
		},
		Timestamp:       time.Unix(int64(parseInt(tx.TimeStamp)), 0),
		Hash:            tx.Hash,
		Nonce:           parseInt(tx.Nonce),
		Index:           parseInt(tx.TransactionIndex),
		From:            tx.From,
		To:              tx.To,
		Value:           parseBig(tx.Value),
		GasLimit:        parseInt(tx.Gas),
		GasUsed:         parseInt(tx.GasUsed),
		GasPrice:        parseBig(tx.GasPrice),
		IsError:         parseBool(tx.IsError),
		Confirmations:   uint64(parseInt(tx.Confirmations)),
		Data:            tx.Input,
		ContractAddress: tx.ContractAddress,
	}
}

func parseTransactionsResponse(r io.Reader) ([]*Transaction, error) {
	res := &transactionsResponse{baseResponse: &baseResponse{}}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}
	transactions := make([]*Transaction, len(res.Transactions))
	for i, tx := range res.Transactions {
		transactions[i] = parseTransaction(tx)
	}
	return transactions, nil
}

func (c *Client) buildTransactionsRequest(addr string, page, offset int) (*http.Request, error) {
	if page <= 0 {
		return nil, errors.New("page param must >= 1")
	}
	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "txlist")
	params.Set("address", addr)
	params.Set("sort", "desc") //newest transactions first
	params.Set("page", fmt.Sprint(page))
	params.Set("offset", fmt.Sprint(offset))
	return c.buildRequest(params)
}

func (c *Client) transactions(ctx context.Context, addr string, page, offset int) ([]*Transaction, error) {
	req, err := c.buildTransactionsRequest(addr, page, offset)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseTransactionsResponse(resp.Body)
}

// Transactions returns a list of transactions to/from the given address
func (c *Client) Transactions(addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(context.Background(), addr, page, offset)
}

// Transactions returns a list of transactions to/from the given address
// with a custom context
func (c *Client) TransactionsContext(ctx context.Context, addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(ctx, addr, page, offset)
}
