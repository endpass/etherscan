package etherscan

// Token represents an ERC20 compatible token
type Token struct {
	Name   string
	Symbol string
	// Number of decimal places used by this token
	Decimals int
}
