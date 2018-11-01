package etherscan

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenSupply(t *testing.T) {
	assert := assert.New(t)

	r := loadTestData(t, "token_totalsupply.json")
	totalSupply, err := parseTokenResponse(r)
	assert.NoError(err)

	val := &big.Int{}
	val.SetString("21265524714464", 10)
	assert.EqualValues(val, totalSupply)
}

func TestTokenBalance(t *testing.T) {
	assert := assert.New(t)

	r := loadTestData(t, "token_totalbalance.json")
	totalBalance, err := parseTokenResponse(r)
	assert.NoError(err)

	val := &big.Int{}
	val.SetString("135499", 10)
	assert.EqualValues(val, totalBalance)
}
