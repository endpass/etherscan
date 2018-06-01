package etherscan

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	assert := assert.New(t)
	data := loadTestData(t, "balance.json")
	expected := &big.Int{}
	expected.SetString("669816163518885498951364", 10)
	bal, err := parseBalanceResponse(data)
	assert.NoError(err)
	assert.EqualValues(expected, bal)
}
