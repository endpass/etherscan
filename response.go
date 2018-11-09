package etherscan

import (
	"errors"
	"math/big"
	"strconv"
)

// Basic fields in all responses
type baseResponse struct {
	Status string `json:"status"`

	Message string `json:"message"`
}

// Checks for error in message field
func checkResponse(resp *baseResponse) error {
	if resp == nil {
		return errors.New("Response is empty")
	}
	if resp.Status != "1" {
		return errors.New("API Error: " + resp.Message)
	}
	return nil
}

// Parse integer and silently discard error
func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// Parse bigInt and silently discard error
func parseBig(s string) *big.Int {
	num := &big.Int{}
	num.SetString(s, 10)
	return num
}

func parseFloat(s string) *big.Float {
	num := &big.Float{}
	num.SetString(s)
	return num
}

func parseIntFromHex(s string) int {
	n, _ := strconv.ParseInt(s, 0, 0)
	return int(n)
}

func parseBigFromHex(s string) *big.Int {
	num := &big.Int{}
	num.SetString(s, 0)
	return num
}

// Parse bool and silently discard error
func parseBool(s string) bool {
	v, _ := strconv.ParseBool(s)
	return v
}
