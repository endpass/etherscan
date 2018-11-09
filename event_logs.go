package etherscan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
)

type eventLogsResponse struct {
	*baseResponse
	Result []eventLog `json:"result"`
}

type eventLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TimeStamp        string   `json:"timeStamp"`
	GasPrice         string   `json:"gasPrice"`
	GasUsed          string   `json:"gasUsed"`
	LogIndex         string   `json:"logIndex"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

type EventLog struct {
	Address          string
	Topics           []string
	Data             string
	BlockNumber      int
	TimeStamp        int
	GasPrice         *big.Int
	GasUsed          int
	LogIndex         int
	TransactionHash  string
	TransactionIndex int
}

func parseEventLogsResponse(r io.Reader) ([]EventLog, error) {
	res := eventLogsResponse{baseResponse: &baseResponse{}}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return nil, err
	}

	if err := checkResponse(res.baseResponse); err != nil {
		return nil, err
	}

	logs := make([]EventLog, len(res.Result))
	for i, l := range res.Result {
		logs[i] = EventLog{
			Address:          l.Address,
			Topics:           l.Topics,
			Data:             l.Data,
			BlockNumber:      parseIntFromHex(l.BlockNumber),
			TimeStamp:        parseIntFromHex(l.TimeStamp),
			GasPrice:         parseBigFromHex(l.GasPrice),
			GasUsed:          parseIntFromHex(l.GasUsed),
			LogIndex:         parseIntFromHex(l.LogIndex),
			TransactionHash:  l.TransactionHash,
			TransactionIndex: parseIntFromHex(l.TransactionIndex),
		}
	}

	return logs, nil
}

type topicOperation string

const (
	TopicOperationAnd topicOperation = "and"
	TopicOperationOr  topicOperation = "or"
)

type EventLogOptions struct {
	FromBlock int
	ToBlock   int
	Address   string

	topics map[string]string
}

func (e *EventLogOptions) addTopic(topic string, op topicOperation) *EventLogOptions {
	if e.topics == nil {
		e.topics = make(map[string]string)
	}

	position := 0
	if len(e.topics) > 0 {
		position = len(e.topics) / 2
	}

	topicKey := fmt.Sprint("topic", position)
	topicKeyOp := fmt.Sprintf("topic%d_%d_opr", position, position+1)

	e.topics[topicKey] = topic
	e.topics[topicKeyOp] = string(op)

	return e
}

func (e *EventLogOptions) AddTopic(topic string) *EventLogOptions {
	return e.addTopic(topic, TopicOperationAnd)
}

func (e *EventLogOptions) AddTopicWithOperation(topic string, op topicOperation) *EventLogOptions {
	return e.addTopic(topic, op)
}

func (c *Client) buildEventLogsRequest(options EventLogOptions) (*http.Request, error) {
	if options.FromBlock == 0 {
		return nil, errors.New("from block required")
	}

	if options.Address == "" && len(options.topics) == 0 {
		return nil, errors.New("address or topics required")
	}

	params := url.Values{}
	params.Set("fromBlock", strconv.Itoa(options.FromBlock))
	params.Set("toBlock", strconv.Itoa(options.ToBlock))
	if options.ToBlock == 0 {
		params.Set("toBlock", "latest")
	}

	if options.Address != "" {
		params.Set("address", options.Address)
	}

	for k, v := range options.topics {
		params.Set(k, v)
	}

	params.Set("module", "logs")
	params.Set("action", "getLogs")

	return c.buildRequest(params)
}

func (c *Client) eventLogs(ctx context.Context, options EventLogOptions) ([]EventLog, error) {
	req, err := c.buildEventLogsRequest(options)
	if err != nil {
		return nil, err
	}
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseEventLogsResponse(resp.Body)
}

// EventLogs returns event logs filtered by options
func (c *Client) EventLogs(options EventLogOptions) ([]EventLog, error) {
	return c.eventLogs(context.Background(), options)
}

// EventLogsContext returns event logs filtered by options with custom context
func (c *Client) EventLogsContext(ctx context.Context, options EventLogOptions) ([]EventLog, error) {
	return c.eventLogs(ctx, options)
}
