package etherscan

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func ExampleClient() {
	client := &Client{
		APIKey: "YOUR-API-KEY",
		// Custom network. Can be mainnet, ropsten, kovan, or rinkeby
		Network: "ropsten",
	}

	address := "0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c"

	balance, err := client.Balance(address)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}
	fmt.Printf("%s Balance: %s", address, balance)

	// Amounts are in big.Int format, check math/big documentation
}

func ExampleClient_context() {
	client := &Client{
		APIKey: "YOUR-API-KEY",
	}

	address := "0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Use context for timeout and cancellation
	balance, err := client.BalanceContext(ctx, address)

	fmt.Print(balance, err)
}
