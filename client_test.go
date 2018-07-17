package etherscan

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Loads fixtures from testdata directory
func loadTestData(t *testing.T, name string) io.Reader {
	path := filepath.Join("testdata", name) // relative path
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

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

func TestBuildRequest(t *testing.T) {
	assert := assert.New(t)
	c := &Client{
		APIKey: "test123",
	}
	params := url.Values{}
	params.Set("module", "account")
	params.Set("action", "balance")
	params.Set("address", "0x123")

	req, err := c.buildRequest(params)
	assert.NoError(err)

	reqURL := req.URL.String()
	assert.Contains(reqURL, apiEndpoints["mainnet"])
	assert.Contains(reqURL, "module=account")
	assert.Contains(reqURL, "action=balance")
	assert.Contains(reqURL, "apikey=test123")
	assert.Contains(reqURL, "address=0x123")
}
