package etherscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultClient(t *testing.T) {
	assert := assert.New(t)
	c := &Client{}
	err := c.setDefaults()
	assert.NoError(err)
	assert.Equal(apiEndpoints["mainnet"], c.apiBase)

	assert.NotNil(c.HTTPClient)
}

func TestClientConfig(t *testing.T) {
	assert := assert.New(t)
	c := &Client{}
	c.Network = "fakenet"
	err := c.setDefaults()
	assert.Error(err)
	assert.Empty(c.apiBase)
}
