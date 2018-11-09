package etherscan

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	assert := assert.New(t)
	r := loadTestData(t, "block.json")

	block, err := parseBlockRewardResponse(r)
	assert.NoError(err)

	assert.Equal(2165403, block.BlockNumber)
	assert.Equal(1472533979, block.TimeStamp)
	assert.Equal("0x13a06d3dfe21e0db5c016c03ea7d2509f7f8d1e3", block.BlockMiner)

	val := &big.Int{}
	val.SetString("5314181600000000000", 10)
	assert.EqualValues(val, block.BlockReward)

	val.SetString("312500000000000000", 10)
	assert.EqualValues(val, block.UncleInclusionReward)

	assert.Len(block.Uncles, 2)

	assert.Equal("0xbcdfc35b86bedf72f0cda046a3c16829a2ef41d1", block.Uncles[0].Miner)
	assert.Equal(0, block.Uncles[0].UnclePosition)
	val.SetString("3750000000000000000", 10)
	assert.EqualValues(val, block.Uncles[0].BlockReward)

	assert.Equal("0x0d0c9855c722ff0c78f21e43aa275a5b8ea60dce", block.Uncles[1].Miner)
	assert.Equal(1, block.Uncles[1].UnclePosition)
	val.SetString("3750000000000000001", 10)
	assert.EqualValues(val, block.Uncles[1].BlockReward)
}
