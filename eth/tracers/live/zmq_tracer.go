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
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/eth/tracers/live/sbe_generated/ethereum_tracing"
	"github.com/ethereum/go-ethereum/log"
)

func init() {
	tracers.LiveDirectory.Register("zmq", newZMQTracer)
}

// zmqTracer is a live tracer that broadcasts events via ZeroMQ
type zmqTracer struct {
	sessionManager *ZMQSessionManager
	eventChan      chan SBEEvent

	// Block processing state (cached from OnBlockStart and OnBlockEndMetrics)
	// Using atomic variables to avoid lock contention
	blockNumber    atomic.Uint64
	blockInsertDur atomic.Int64 // Nanoseconds
}

// zmqTracerConfig holds tracer configuration
type zmqTracerConfig struct {
	BindEndpoints []string `json:"bindEndpoints"` // ZMQ ROUTER endpoints
	QueueSize     int      `json:"queueSize"`     // Event queue size
	MaxSessions   int      `json:"maxSessions"`   // Maximum concurrent sessions
	SendTimeout   int      `json:"sendTimeout"`   // Send timeout in ms (default 2000ms)
}

// newZMQTracer creates a new ZMQ tracer instance
func newZMQTracer(cfg json.RawMessage) (*tracing.Hooks, error) {
	var config zmqTracerConfig
	if err := json.Unmarshal(cfg, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	// Set defaults
	if len(config.BindEndpoints) == 0 {
		config.BindEndpoints = []string{"ipc:///tmp/geth.tracer.sock"}
	}
	if config.QueueSize == 0 {
		config.QueueSize = 1000
	}
	if config.MaxSessions == 0 {
		config.MaxSessions = 100
	}

	// Create session manager
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: config.BindEndpoints,
		QueueSize:     config.QueueSize,
		MaxSessions:   config.MaxSessions,
		SendTimeout:   config.SendTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session manager: %v", err)
	}

	// Start ZMQ session manager
	if err := sessionMgr.Start(); err != nil {
		return nil, fmt.Errorf("failed to start session manager: %v", err)
	}

	tracer := &zmqTracer{
		sessionManager: sessionMgr,
		eventChan:      sessionMgr.eventChan,
	}

	log.Info("ZMQ tracer initialized",
		"bindEndpoints", config.BindEndpoints,
		"queueSize", config.QueueSize,
		"maxSessions", config.MaxSessions)

	return &tracing.Hooks{
		OnTxStart:         tracer.OnTxStart,
		OnTxEnd:           tracer.OnTxEnd,
		OnEnter:           tracer.OnEnter,
		OnExit:            tracer.OnExit,
		OnBlockStart:      tracer.OnBlockStart,
		OnBlockEnd:        tracer.OnBlockEnd,
		OnBlockEndMetrics: tracer.OnBlockEndMetrics,
		OnLog:             tracer.OnLog,
		OnClose:           tracer.OnClose,
	}, nil
}

// OnTxStart is called when a transaction starts executing
func (t *zmqTracer) OnTxStart(vm *tracing.VMContext, tx *types.Transaction, from common.Address) {
	// Fast path: skip if no sessions subscribe to TxStart events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.TxStart) {
		return
	}

	// Convert common types to [N]uint8 arrays
	var txHash [32]uint8
	copy(txHash[:], tx.Hash().Bytes())

	var fromAddr [20]uint8
	copy(fromAddr[:], from.Bytes())

	event := &ethereum_tracing.TxStartEvent{
		EventType: ethereum_tracing.EventType.TxStart,
		Timestamp: uint64(time.Now().UnixNano()),
		TxHash:    txHash,
		From:      fromAddr,
		GasLimit:  tx.Gas(),
		Nonce:     tx.Nonce(),
	}

	// Handle optional fields
	if to := tx.To(); to != nil {
		copy(event.To[:], to.Bytes())
	}
	if value := tx.Value(); value != nil {
		copy(event.Value[:], value.Bytes())
	}
	if gasPrice := tx.GasPrice(); gasPrice != nil {
		copy(event.GasPrice[:], gasPrice.Bytes())
	}

	t.sendEvent(event)
}

// OnTxEnd is called when a transaction completes
func (t *zmqTracer) OnTxEnd(receipt *types.Receipt, err error) {
	// Fast path: skip if no sessions subscribe to TxEnd events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.TxEnd) {
		return
	}

	status := uint8(0)
	if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
		status = 1
	}

	var txHash [32]uint8
	var blockHash [32]uint8

	if receipt != nil {
		copy(txHash[:], receipt.TxHash.Bytes())
		if receipt.BlockHash != (common.Hash{}) {
			copy(blockHash[:], receipt.BlockHash.Bytes())
		}
	}

	event := &ethereum_tracing.TxEndEvent{
		EventType:   ethereum_tracing.EventType.TxEnd,
		Timestamp:   uint64(time.Now().UnixNano()),
		TxHash:      txHash,
		Status:      status,
		GasUsed:     receipt.GasUsed,
		BlockNumber: receipt.BlockNumber.Uint64(),
		BlockHash:   blockHash,
	}

	t.sendEvent(event)
}

// OnEnter is called when entering a call frame
func (t *zmqTracer) OnEnter(depth int, typ byte, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	// Fast path: skip if no sessions subscribe to Enter events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.Enter) {
		return
	}

	var fromAddr [20]uint8
	copy(fromAddr[:], from.Bytes())

	var toAddr [20]uint8
	copy(toAddr[:], to.Bytes())

	event := &ethereum_tracing.CallEvent{
		EventType: ethereum_tracing.EventType.Enter,
		Timestamp: uint64(time.Now().UnixNano()),
		Depth:     uint8(depth),
		CallType:  typ,
		From:      fromAddr,
		To:        toAddr,
		Gas:       gas,
		Input:     input,
	}

	if value != nil {
		copy(event.Value[:], value.Bytes())
	}

	t.sendEvent(event)
}

// OnExit is called when exiting a call frame
func (t *zmqTracer) OnExit(depth int, output []byte, gasUsed uint64, err error, reverted bool) {
	// Fast path: skip if no sessions subscribe to Exit events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.Exit) {
		return
	}

	var revertedFlag uint8
	if reverted {
		revertedFlag = 1
	}

	event := &ethereum_tracing.CallEvent{
		EventType:  ethereum_tracing.EventType.Exit,
		Timestamp:  uint64(time.Now().UnixNano()),
		Depth:      uint8(depth),
		OutputSize: uint32(len(output)),
		Reverted:   revertedFlag,
		Output:     output,
	}

	t.sendEvent(event)
}

// OnBlockStart is called when block processing starts
func (t *zmqTracer) OnBlockStart(event tracing.BlockEvent) {
	b := event.Block

	// Cache block number for OnBlockEndMetrics and OnBlockEnd
	t.blockNumber.Store(b.NumberU64())

	// Fast path: skip if no sessions subscribe to BlockStart events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.BlockStart) {
		return
	}

	var hash [32]uint8
	copy(hash[:], b.Hash().Bytes())

	var parentHash [32]uint8
	copy(parentHash[:], b.ParentHash().Bytes())

	sbeEvent := &ethereum_tracing.BlockStartEvent{
		EventType:      ethereum_tracing.EventType.BlockStart,
		Timestamp:      uint64(time.Now().UnixNano()),
		Number:         b.NumberU64(),
		Hash:           hash,
		ParentHash:     parentHash,
		BlockTimestamp: b.Time(),
		TxCount:        uint16(b.Transactions().Len()),
	}

	t.sendEvent(sbeEvent)
}

// OnBlockEnd is called when block processing ends
// This is the last hook to execute (defer), responsible for sending BlockEndEvent
func (t *zmqTracer) OnBlockEnd(err error) {
	// Fast path: skip if no sessions subscribe to BlockEnd events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.BlockEnd) {
		return
	}

	// Read cached values (clear after reading to avoid stale data)
	blockNumber := t.blockNumber.Load()
	insertDur := time.Duration(t.blockInsertDur.Load())

	// Clear cached values after reading
	t.blockNumber.Store(0)
	t.blockInsertDur.Store(0)

	// Prepare error message
	var errorMsg []uint8
	if err != nil {
		errorMsg = []uint8(err.Error())
	}

	sbeEvent := &ethereum_tracing.BlockEndEvent{
		EventType:        ethereum_tracing.EventType.BlockEnd,
		Timestamp:        uint64(time.Now().UnixNano()),
		Number:           blockNumber,
		InsertDurationNs: uint64(insertDur),
		ErrorMsg:         errorMsg,
	}

	t.sendEvent(sbeEvent)
}

// OnBlockEndMetrics is called after block processing with metrics data
// Caches the duration for OnBlockEnd to use
func (t *zmqTracer) OnBlockEndMetrics(blockNumber uint64, blockInsertDuration time.Duration) {
	// Cache duration for OnBlockEnd to use
	t.blockInsertDur.Store(int64(blockInsertDuration))
}

// OnLog is called when a log is emitted
func (t *zmqTracer) OnLog(log *types.Log) {
	// Fast path: skip if no sessions subscribe to Log events
	if !t.sessionManager.HasSessionsForEventType(ethereum_tracing.EventType.Log) {
		return
	}

	var address [20]uint8
	copy(address[:], log.Address.Bytes())

	// Convert topics
	topicsCount := uint8(len(log.Topics))
	topics := make([]ethereum_tracing.LogEventTopics, topicsCount)
	for i, topic := range log.Topics {
		copy(topics[i].Topic[:], topic.Bytes())
	}

	event := &ethereum_tracing.LogEvent{
		EventType:   ethereum_tracing.EventType.Log,
		Timestamp:   uint64(time.Now().UnixNano()),
		Address:     address,
		TopicsCount: topicsCount,
		Topics:      topics,
		LogData:     log.Data,
	}

	t.sendEvent(event)
}

// OnClose is called when the tracer is closed
func (t *zmqTracer) OnClose() {
	t.sessionManager.Stop()
}

// sendEvent sends an event to the event channel (non-blocking)
func (t *zmqTracer) sendEvent(event SBEEvent) {
	select {
	case t.eventChan <- event:
	default:
		log.Warn("Event channel full, dropping event")
	}
}
