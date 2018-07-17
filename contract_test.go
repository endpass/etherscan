package etherscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractABI(t *testing.T) {
	assert := assert.New(t)
	r := loadTestData(t, "abi.json")
	rawABI, err := parseABIResponse(r)
	assert.NoError(err)
	assert.NotEmpty(rawABI)
	assert.Contains(string(rawABI), "inputs")
}
