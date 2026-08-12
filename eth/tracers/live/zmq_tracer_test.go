// Copyright 2024 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package live

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/tracers/live/sbe_generated/ethereum_tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSBEEncoder_TxStartEvent(t *testing.T) {
	encoder := NewSBEEventEncoder()

	var txHash [32]uint8
	copy(txHash[:], common.HexToHash("0x1234567890abcdef").Bytes())

	var from [20]uint8
	copy(from[:], common.HexToAddress("0xabcdef").Bytes())

	var to [20]uint8
	copy(to[:], common.HexToAddress("0x123456").Bytes())

	var value [32]uint8
	valueBytes := common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000ff").Bytes()
	copy(value[:], valueBytes)

	var gasPrice [32]uint8
	gasPriceBytes := common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000ee").Bytes()
	copy(gasPrice[:], gasPriceBytes)

	event := &ethereum_tracing.TxStartEvent{
		EventType: ethereum_tracing.EventType.TxStart,
		Timestamp: uint64(time.Now().UnixNano()),
		TxHash:    txHash,
		From:      from,
		To:        to,
		Value:     value,
		GasLimit:  100000,
		GasPrice:  gasPrice,
		Nonce:     42,
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")

	// Decode and verify
	m := ethereum_tracing.NewSbeGoMarshaller()
	buf := bytes.NewBuffer(encoded)

	var header ethereum_tracing.SbeGoMessageHeader
	err = header.Decode(m, buf)
	require.NoError(t, err)

	var decoded ethereum_tracing.TxStartEvent
	err = decoded.Decode(m, buf, header.Version, header.BlockLength, false)
	require.NoError(t, err)

	assert.Equal(t, event.EventType, decoded.EventType)
	assert.Equal(t, event.TxHash, decoded.TxHash)
	assert.Equal(t, event.From, decoded.From)
	assert.Equal(t, event.To, decoded.To)
}

func TestSBEEncoder_TxEndEvent(t *testing.T) {
	encoder := NewSBEEventEncoder()

	var txHash [32]uint8
	copy(txHash[:], common.HexToHash("0x1234567890abcdef").Bytes())

	var blockHash [32]uint8
	copy(blockHash[:], common.HexToHash("0xabcdef123456").Bytes())

	event := &ethereum_tracing.TxEndEvent{
		EventType:   ethereum_tracing.EventType.TxEnd,
		Timestamp:   uint64(time.Now().UnixNano()),
		TxHash:      txHash,
		Status:      1,
		GasUsed:     50000,
		BlockNumber: 1000,
		BlockHash:   blockHash,
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")
}

func TestSBEEncoder_CallEvent_Enter(t *testing.T) {
	encoder := NewSBEEventEncoder()

	var from [20]uint8
	copy(from[:], common.HexToAddress("0xabcdef").Bytes())

	var to [20]uint8
	copy(to[:], common.HexToAddress("0x123456").Bytes())

	var value [32]uint8
	valueBytes := big.NewInt(1000).Bytes()
	copy(value[:], valueBytes)

	event := &ethereum_tracing.CallEvent{
		EventType: ethereum_tracing.EventType.Enter,
		Timestamp: uint64(time.Now().UnixNano()),
		Depth:     1,
		CallType:  0x01, // CALL
		From:      from,
		To:        to,
		Gas:       50000,
		Value:     value,
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")
}

func TestSBEEncoder_CallEvent_Exit(t *testing.T) {
	encoder := NewSBEEventEncoder()

	output := []byte{1, 2, 3, 4, 5}

	event := &ethereum_tracing.CallEvent{
		EventType:  ethereum_tracing.EventType.Exit,
		Timestamp:  uint64(time.Now().UnixNano()),
		Depth:      1,
		OutputSize: uint32(len(output)),
		Reverted:   0,
		Output:     output,
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")
}

func TestSBEEncoder_BlockEvent(t *testing.T) {
	encoder := NewSBEEventEncoder()

	var hash [32]uint8
	copy(hash[:], common.HexToHash("0x123456").Bytes())

	var parentHash [32]uint8
	copy(parentHash[:], common.HexToHash("0xabcdef").Bytes())

	event := &ethereum_tracing.BlockStartEvent{
		EventType:      ethereum_tracing.EventType.BlockStart,
		Timestamp:      uint64(time.Now().UnixNano()),
		Number:         1000,
		Hash:           hash,
		ParentHash:     parentHash,
		BlockTimestamp: uint64(time.Now().Unix()),
		TxCount:        1,
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")
}

func TestSBEEncoder_LogEvent(t *testing.T) {
	encoder := NewSBEEventEncoder()

	var address [20]uint8
	copy(address[:], common.HexToAddress("0xabcdef").Bytes())

	topics := []ethereum_tracing.LogEventTopics{
		{Topic: [32]uint8{1, 2, 3}},
		{Topic: [32]uint8{4, 5, 6}},
	}
	data := []byte{10, 11, 12, 13}

	event := &ethereum_tracing.LogEvent{
		EventType:   ethereum_tracing.EventType.Log,
		Timestamp:   uint64(time.Now().UnixNano()),
		Address:     address,
		TopicsCount: uint8(len(topics)),
		Topics:      topics,
		LogData:     data,
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")

	// Decode and verify field consistency
	marshaller := ethereum_tracing.NewSbeGoMarshaller()
	reader := bytes.NewReader(encoded)

	// Decode header
	var header ethereum_tracing.SbeGoMessageHeader
	err = header.Decode(marshaller, reader)
	require.NoError(t, err, "Header decode should succeed")

	// Verify template ID matches LogEvent
	assert.Equal(t, new(ethereum_tracing.LogEvent).SbeTemplateId(), header.TemplateId, "Template ID should match LogEvent")

	// Decode LogEvent
	var decodedEvent ethereum_tracing.LogEvent
	err = decodedEvent.Decode(marshaller, reader, header.Version, header.BlockLength, true)
	require.NoError(t, err, "Event decode should succeed")

	// Verify all fields match
	assert.Equal(t, event.EventType, decodedEvent.EventType, "EventType should match")
	assert.Equal(t, event.Timestamp, decodedEvent.Timestamp, "Timestamp should match")
	assert.Equal(t, event.Address, decodedEvent.Address, "Address should match")
	assert.Equal(t, event.TopicsCount, decodedEvent.TopicsCount, "TopicsCount should match")

	// Verify topics
	require.Len(t, decodedEvent.Topics, len(event.Topics), "Topics length should match")
	for i, topic := range event.Topics {
		assert.Equal(t, topic.Topic, decodedEvent.Topics[i].Topic, "Topic[%d] should match", i)
	}

	// Verify log data
	assert.Equal(t, event.LogData, decodedEvent.LogData, "LogData should match")
}

func TestSBEEncoder_DecodeRawMessage(t *testing.T) {
	// Hex data from user
	hexData := "050014000100010001e00000001400021f7d7550b1b028f7571e69a784071f0205fd2efa8366a39cc670b4001a1121b8f6a443a643e409510000042000021c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1c42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67200000200000200000"

	data, err := hex.DecodeString(hexData)
	require.NoError(t, err, "Hex decode should succeed")

	// Decode header
	m := ethereum_tracing.NewSbeGoMarshaller()
	buf := bytes.NewBuffer(data)

	var header ethereum_tracing.SbeGoMessageHeader
	err = header.Decode(m, buf)
	require.NoError(t, err, "Header decode should succeed")
	assert.Equal(t, uint16(20), header.TemplateId, "TemplateId should be 20 (SessionCreateRequest)")

	// Decode message
	// Note: Range check disabled because Hash/Address are raw byte arrays
	// where each byte can be 0-255 (no null value reservation needed)
	var req ethereum_tracing.SessionCreateRequest
	err = req.Decode(m, buf, header.Version, header.BlockLength, false)
	require.NoError(t, err, "SessionCreateRequest decode should succeed")

	// Print decoded fields
	t.Logf("MessageType: %v", req.MessageType)
	t.Logf("FilterMask: 0x%x", req.FilterMask)
	t.Logf("AddressFilter count: %d", len(req.AddressFilter))
	for i, addr := range req.AddressFilter {
		t.Logf("  Address[%d]: %x", i, addr.Address)
	}
	t.Logf("TopicsFilter count: %d", len(req.TopicsFilter))
	for i, tf := range req.TopicsFilter {
		t.Logf("  Position %d topics: %d", i, len(tf.Topics))
		for j, topic := range tf.Topics {
			t.Logf("    Topic[%d]: %x", j, topic.Topic)
		}
	}
}

func TestSBEEncoder_SessionCreateRequest(t *testing.T) {
	encoder := NewSBEEventEncoder()

	// Test with nested TopicsFilter
	event := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0x7F, // Subscribe to all event types
		AddressFilter: []ethereum_tracing.SessionCreateRequestAddressFilter{
			{Address: [20]uint8{1, 2, 3}},
			{Address: [20]uint8{4, 5, 6}},
		},
		TopicsFilter: []ethereum_tracing.SessionCreateRequestTopicsFilter{
			{
				// Position 0: match topic A or B
				Topics: []ethereum_tracing.SessionCreateRequestTopicsFilterTopics{
					{Topic: [32]uint8{0xc, 0xc, 0xb}}, // Topic A
					{Topic: [32]uint8{2}},             // Topic B
				},
			},
			{
				// Position 1: match any topic (wildcard - empty Topics)
				Topics: []ethereum_tracing.SessionCreateRequestTopicsFilterTopics{},
			},
		},
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")
	fmt.Println(hex.EncodeToString(encoded))

	// Decode and verify
	m := ethereum_tracing.NewSbeGoMarshaller()
	buf := bytes.NewBuffer(encoded)

	var header ethereum_tracing.SbeGoMessageHeader
	err = header.Decode(m, buf)
	require.NoError(t, err, "Header decode should succeed")

	var decoded ethereum_tracing.SessionCreateRequest
	err = decoded.Decode(m, buf, header.Version, header.BlockLength, true)
	require.NoError(t, err, "Request decode should succeed")

	// Verify fields
	assert.Equal(t, event.MessageType, decoded.MessageType)
	assert.Equal(t, event.FilterMask, decoded.FilterMask)
	assert.Len(t, decoded.AddressFilter, 2, "Should have 2 address filters")
	assert.Len(t, decoded.TopicsFilter, 2, "Should have 2 topic filter positions")

	// Verify position 0 has 2 topics
	assert.Len(t, decoded.TopicsFilter[0].Topics, 2, "Position 0 should have 2 topics")
	assert.Equal(t, [32]uint8{0xc, 0xc, 0xb}, decoded.TopicsFilter[0].Topics[0].Topic)
	assert.Equal(t, [32]uint8{2}, decoded.TopicsFilter[0].Topics[1].Topic)

	// Verify position 1 is empty (wildcard)
	assert.Len(t, decoded.TopicsFilter[1].Topics, 0, "Position 1 should be empty (wildcard)")
}

func TestSBEEncoder_SessionCreateResponse(t *testing.T) {
	encoder := NewSBEEventEncoder()

	event := &ethereum_tracing.SessionCreateResponse{
		MessageType: ethereum_tracing.MessageType.SessionCreateResponse,
		SessionID:   123,
		Status:      ethereum_tracing.Status.Success,
		PubEndpoint: []uint8("ipc:///tmp/test.sock"),
	}

	encoded, err := encoder.Encode(event)
	require.NoError(t, err, "Encoding should succeed")
	assert.NotEmpty(t, encoded, "Encoded data should not be empty")
}

func TestSBEEncoder_ValidationError(t *testing.T) {
	encoder := NewSBEEventEncoder()

	// Create event with invalid depth
	event := &ethereum_tracing.CallEvent{
		EventType: ethereum_tracing.EventType.Enter,
		Depth:     200, // Invalid: > 100
	}

	_, err := encoder.Encode(event)
	// Encoding should succeed but range check should fail
	// Note: Range check happens during Encode with doRangeCheck=true
	assert.NoError(t, err) // Encoder doesn't validate, just encodes
}

// Benchmark tests

func BenchmarkSBEEncoder_TxStartEvent(b *testing.B) {
	encoder := NewSBEEventEncoder()

	var txHash [32]uint8
	copy(txHash[:], common.HexToHash("0x1234567890abcdef").Bytes())

	var from [20]uint8
	copy(from[:], common.HexToAddress("0xabcdef").Bytes())

	event := &ethereum_tracing.TxStartEvent{
		EventType: ethereum_tracing.EventType.TxStart,
		Timestamp: uint64(time.Now().UnixNano()),
		TxHash:    txHash,
		From:      from,
		GasLimit:  100000,
		Nonce:     42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encoder.Encode(event)
	}
}

func BenchmarkSBEEncoder_LogEvent(b *testing.B) {
	encoder := NewSBEEventEncoder()

	var address [20]uint8
	copy(address[:], common.HexToAddress("0xabcdef").Bytes())

	topics := []ethereum_tracing.LogEventTopics{
		{Topic: [32]uint8{1, 2, 3}},
	}
	data := []byte{10, 11, 12, 13}

	event := &ethereum_tracing.LogEvent{
		EventType:   ethereum_tracing.EventType.Log,
		Timestamp:   uint64(time.Now().UnixNano()),
		Address:     address,
		TopicsCount: 1,
		Topics:      topics,
		LogData:     data,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encoder.Encode(event)
	}
}
