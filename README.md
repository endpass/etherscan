# etherscan

[![Build Status](https://travis-ci.org/endpass/etherscan.svg?branch=master)](https://travis-ci.org/endpass/etherscan)

Go client library for the Etherscan.io Ethereum block explorer API.
This library aims to be full-featured, well-tested, and easy to use.

[**Documentation**](https://godoc.org/github.com/endpass/etherscan).

## Install
`go get github.com/endpass/etherscan`

## Usage

```go
import (
    "fmt"
    "github.com/endpass/etherscan"
)

func main() {
	client := &Client{
		APIKey: "YOUR-API-KEY",
		// Custom network. Can be mainnet, ropsten, kovan, or rinkeby
    // Defaults to mainnet if not set
		Network: "ropsten",
	}

	address := "0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c"

	balance, err := client.Balance(address)
  if err != nil {
      fmt.Printf("ERROR: %s", err)
      return
  }
  fmt.Printf("%s Balance: %s", address, balance)
}
```

All client commands have a ..Context version that allows their use with
context.Context.

## Status
Supported featues of the [Etherscan API](https://etherscan.io/apis):
- [x] Accounts
- [x] Contracts
- [x] Transactions
- [ ] Blocks
- [ ] Event Logs
- [ ] Tokens
- [ ] Stats
