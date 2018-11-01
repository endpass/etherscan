package etherscan

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsTotalSupply(t *testing.T) {
	assert := assert.New(t)

	r := loadTestData(t, "stats_totalsupply.json")
	total, err := parseStatsTotalSupplyResponse(r)
	assert.NoError(err)

	val := &big.Int{}
	val.SetString("102935195936600000000000000", 10)

	assert.EqualValues(val, total)
}

func TestStatsLastPrice(t *testing.T) {
	assert := assert.New(t)

	r := loadTestData(t, "stats_lastprice.json")
	lastPrice, err := parseStatsLastPriceResponse(r)
	assert.NoError(err)

	val := &big.Float{}
	val.SetString("0.03118")
	assert.EqualValues(val, lastPrice.Ethbtc)
	assert.Equal(1541092070, lastPrice.EthbtcTimestamp)

	val.SetString("197.48")
	assert.EqualValues(val, lastPrice.Ethusd)
	assert.Equal(1541092064, lastPrice.EthusdTimestamp)
}
