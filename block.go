package etherscan

// Block is a single block in the chain
type Block struct {
	Number int
	Hash   string
	// All transactions mined in this block
	Transactions []*Transaction
}
