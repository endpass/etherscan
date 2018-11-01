package etherscan

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventLogs(t *testing.T) {
	assert := assert.New(t)

	r := loadTestData(t, "event_logs.json")
	logs, err := parseEventLogsResponse(r)
	assert.NoError(err)

	assert.Len(logs, 369)
	log := logs[0]

	val := &big.Int{}
	val.SetString("50000000000", 10)
	assert.EqualValues(val, log.GasPrice)

	assert.Equal("0x33990122638b9132ca29c723bdf037f1a891a70c", log.Address)
	assert.Len(log.Topics, 3)
	assert.Equal("0xf63780e752c6a54a94fc52715dbc5518a3b4c3c2833d301a204226548a2a8545", log.Topics[0])
	assert.Equal("0x72657075746174696f6e00000000000000000000000000000000000000000000", log.Topics[1])
	assert.Equal("0x000000000000000000000000d9b2f59f3b5c7b3c67047d2f03c3e8052470be92", log.Topics[2])
	assert.Equal("0x", log.Data)
	assert.Equal(379224, log.BlockNumber)
	assert.Equal(1444767884, log.TimeStamp)
	assert.Equal(67202, log.GasUsed)
	assert.Equal(0, log.LogIndex)
	assert.Equal("0x0b03498648ae2da924f961dda00dc6bb0a8df15519262b7e012b7d67f4bb7e83", log.TransactionHash)
	assert.Equal(0, log.TransactionIndex)
}

func TestEventLogOptions(t *testing.T) {
	assert := assert.New(t)
	options := &EventLogOptions{}

	options.AddTopic("1").
		AddTopicWithOperation("2", TopicOperationOr).
		AddTopic("3")

	assert.Equal(6, len(options.topics))

	checks := map[string]string{
		"topic0":       "1",
		"topic0_1_opr": "and",
		"topic1":       "2",
		"topic1_2_opr": "or",
		"topic2":       "3",
		"topic2_3_opr": "and",
	}

	for k, v := range checks {
		optionsValue, ok := options.topics[k]
		if !ok {
			t.Fatalf("key %s not in topics", k)
		}
		assert.Equal(v, optionsValue)
	}
}
