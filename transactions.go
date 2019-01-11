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

// Types of transactions indexed by Etherscan
type txType uint8

const (
	txNormal = iota
	txInternal
	txToken
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
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	ErrCode           string `json:"errCode"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	Type              string `json:"type"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	TraceID           string `json:"traceId"`
	Confirmations     string `json:"confirmations"`
}

// Transaction represents a standard Ethereum transaction
type Transaction struct {
	Block *Block

	Token *Token

	Internal *InternalTransaction

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

	// Detailed error from contract
	Error error

	Confirmations uint64

	// Data sent to transaction, encoded as hex
	Data string
}

// Internal transaction is a value transfer inside a contract's code
type InternalTransaction struct {
	// Transaction type, such as "call" for a method call
	Type    string
	TraceID string
}

func parseTransaction(tx *transactionResponse) *Transaction {
	parsedTx := &Transaction{
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
	// Transaction only has a block if it is confirmed
	if tx.BlockNumber != "" {
		parsedTx.Block = &Block{
			Number: parseInt(tx.BlockNumber),
			Hash:   tx.BlockHash,
		}
	}
	// ERC20 token transactions, in this case Value is the amount of tokens
	// transfered
	if tx.TokenSymbol != "" {
		parsedTx.Token = &Token{
			Name:     tx.TokenName,
			Symbol:   tx.TokenSymbol,
			Decimals: parseInt(tx.TokenDecimal),
		}
	}
	// Internal transactions should always have a Type
	if tx.Type != "" {
		parsedTx.Internal = &InternalTransaction{
			Type:    tx.Type,
			TraceID: tx.TraceID,
		}
	}

	if tx.ErrCode != "" {
		parsedTx.Error = errors.New(tx.ErrCode)
	}
	return parsedTx
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

func (c *Client) buildTransactionsRequest(addr string, page, offset int, category txType) (*http.Request, error) {
	var action string
	if page <= 0 {
		return nil, errors.New("page param must >= 1")
	}
	switch category {
	case txNormal:
		action = "txlist"
	case txInternal:
		action = "txlistinternal"
	case txToken:
		action = "tokentx"
	}
	if action == "" {
		return nil, errors.New("Unsupported transaction category")
	}
	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", action)
	params.Set("address", addr)
	params.Set("sort", "desc") //newest transactions first
	params.Set("page", fmt.Sprint(page))
	params.Set("offset", fmt.Sprint(offset))
	return c.buildRequest(params)
}

func (c *Client) transactions(ctx context.Context, addr string, page, offset int, category txType) ([]*Transaction, error) {
	req, err := c.buildTransactionsRequest(addr, page, offset, category)
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

// Transactions returns a list of standard transactions to/from the given address
func (c *Client) Transactions(addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(context.Background(), addr, page, offset, txNormal)
}

// TransactionsContext returns a list of standard transactions to/from the given address
// with a custom context
func (c *Client) TransactionsContext(ctx context.Context, addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(ctx, addr, page, offset, txNormal)
}

// TokenTransactions returns a list of ERC20 token transactions to/from the given address
func (c *Client) TokenTransactions(addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(context.Background(), addr, page, offset, txToken)
}

// TokenTransactionsContext returns a list of ERC20 token transactions to/from the given address
// with a custom context
func (c *Client) TokenTransactionsContext(ctx context.Context, addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(ctx, addr, page, offset, txToken)
}

// InternalTransactions returns a list of internal contract transactions for
// the contract at the given address
func (c *Client) InternalTransactions(addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(context.Background(), addr, page, offset, txInternal)
}

// InternalTransactionsContext returns a list of internal contract transactions for
// the contract at the given address
// with a custom context
func (c *Client) InternalTransactionsContext(ctx context.Context, addr string, page, offset int) ([]*Transaction, error) {
	return c.transactions(ctx, addr, page, offset, txInternal)
}
