package etherscan

import "errors"

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
