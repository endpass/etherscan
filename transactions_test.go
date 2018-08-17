package etherscan

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactions(t *testing.T) {
	assert := assert.New(t)
	r := loadTestData(t, "transactions.json")

	txs, err := parseTransactionsResponse(r)
	assert.NoError(err)
	assert.Len(txs, 3)

	tx := txs[0]
	val := &big.Int{}
	val.SetString("10000000000000000000000", 10)

	assert.Equal(1959393, tx.Block.Number)
	assert.Equal("0x03dd2d32f8ea317eab05180152b6e959751a32f3de202b3d7caf3b748273e40a", tx.Block.Hash)
	assert.EqualValues(time.Unix(1469591193, 0), tx.Timestamp)
	assert.Equal("0x81b904768ecbba7a8ddd69bec5c8bb63e5b28bb973caae4010e4e56c55fc0462", tx.Hash)
	assert.Equal(19, tx.Nonce)
	assert.Equal(3, tx.Index)
	assert.Equal("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a", tx.From)
	assert.Equal("0x1bb0ac60363e320bc45fdb15aed226fb59c88e44", tx.To)
	assert.EqualValues(val, tx.Value)
	assert.Equal(127964, tx.GasLimit)
	assert.Equal(27964, tx.GasUsed)
	assert.EqualValues(big.NewInt(20000000000), tx.GasPrice)
	assert.Equal(false, tx.IsError)
	assert.EqualValues(4018297, tx.Confirmations)
}

func ExampleClient_Transactions() {
	client := &Client{
		APIKey: "YOUR-API-KEY",
	}

	address := "0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c"

	page := 1
	limit := 10

	transactions, err := client.Transactions(address, page, limit)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	if len(transactions) == 0 {
		return
	}

	tx := transactions[0]
	fmt.Printf(`
	Transaction Hash: %s
	Block Number: %d
	Time: %s
	From: %s
	To: %s
	Value: %s
	Confirmations: %d

	`, tx.Hash, tx.Block.Number, tx.Timestamp, tx.From, tx.To, tx.Value, tx.Confirmations)
}
