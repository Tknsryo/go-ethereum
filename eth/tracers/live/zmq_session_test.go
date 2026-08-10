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
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth/tracers/live/sbe_generated/ethereum_tracing"
	zmq "gopkg.in/pebbe/zmq4.v1"
)

// TestZMQTracer_ClientDisconnect 测试客户端断开连接场景
// 流程：
// 1. 客户端连接并创建会话
// 2. 服务端推送 BlockEvent
// 3. 客户端接收事件
// 4. 客户端断开连接
// 5. 服务端检测到断开并清理会话
func TestZMQTracer_ClientDisconnect(t *testing.T) {
	// 1. 创建 session manager
	endpoint := fmt.Sprintf("ipc:///tmp/test_tracer_%d.sock", time.Now().UnixNano())

	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{endpoint},
		QueueSize:     100,
		MaxSessions:   10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// 启动 session manager
	if err := sessionMgr.Start(); err != nil {
		t.Fatalf("Failed to start session manager: %v", err)
	}
	defer sessionMgr.Stop()

	// 等待 socket 绑定完成
	time.Sleep(100 * time.Millisecond)

	// 2. 客户端连接
	client, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		t.Fatalf("Failed to create client socket: %v", err)
	}

	// 设置唯一的 identity
	clientIdentity := fmt.Sprintf("test-client-%d", time.Now().UnixNano())
	if err := client.SetIdentity(clientIdentity); err != nil {
		t.Fatalf("Failed to set client identity: %v", err)
	}

	if err := client.Connect(endpoint); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// 等待连接建立
	time.Sleep(100 * time.Millisecond)

	// 3. 发送会话创建请求
	// DEALER 发送：自动添加 delimiter
	t.Log("Client: Sending session create request")
	request := createTestSessionRequest()
	if _, err := client.SendBytes(request, 0); err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// 4. 接收会话创建响应
	// DEALER 接收：自动去掉 identity 和 delimiter，只接收 data
	t.Log("Client: Waiting for session create response")
	client.SetRcvtimeo(5 * time.Second)

	response, err := client.RecvBytes(0)
	if err != nil {
		t.Fatalf("Failed to receive response: %v", err)
	}

	sessionID, err := parseSessionResponse(response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	t.Logf("Client: Session created with ID %d", sessionID)

	// 验证 session 创建成功
	if sessionMgr.SessionCount() != 1 {
		t.Errorf("Expected 1 session, got %d", sessionMgr.SessionCount())
	}

	// 5. 服务端推送 BlockEvent
	t.Log("Server: Sending block event")
	blockEvent := createTestBlockEvent()
	sessionMgr.eventChan <- blockEvent

	// 6. 客户端接收 BlockEvent
	t.Log("Client: Waiting for block event")
	// 设置超时，防止无限等待
	client.SetRcvtimeo(2 * time.Second)

	eventData, err := client.RecvBytes(0)
	if err != nil {
		t.Fatalf("Failed to receive event: %v", err)
	}

	receivedEvent, err := parseBlockEvent(eventData)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}
	t.Logf("Client: Received block event, block number=%d", receivedEvent.Number)

	// 验证事件内容
	if receivedEvent.Number != blockEvent.Number {
		t.Errorf("Block number mismatch: got %d, want %d", receivedEvent.Number, blockEvent.Number)
	}

	// 7. 客户端断开连接
	t.Log("Client: Disconnecting")
	client.Close()

	// 等待断开检测
	time.Sleep(200 * time.Millisecond)

	// 8. 服务端推送另一个事件，应该触发断线检测
	t.Log("Server: Sending another event (should trigger disconnect detection)")
	anotherEvent := createTestBlockEvent()
	anotherEvent.Number = 999
	sessionMgr.eventChan <- anotherEvent

	// 等待事件处理和会话清理
	time.Sleep(300 * time.Millisecond)

	// 9. 验证会话已清理
	if sessionMgr.SessionCount() != 0 {
		t.Errorf("Expected 0 sessions after disconnect, got %d", sessionMgr.SessionCount())
	}
	t.Log("Server: Session cleaned up successfully")

	// 10. 验证 IsEmpty
	if !sessionMgr.IsEmpty() {
		t.Error("Expected IsEmpty() to return true")
	}
}

// TestZMQTracer_MultipleClients 测试多个客户端场景
func TestZMQTracer_MultipleClients(t *testing.T) {
	endpoint := fmt.Sprintf("ipc:///tmp/test_tracer_multi_%d.sock", time.Now().UnixNano())

	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{endpoint},
		QueueSize:     100,
		MaxSessions:   10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	if err := sessionMgr.Start(); err != nil {
		t.Fatalf("Failed to start session manager: %v", err)
	}
	defer sessionMgr.Stop()

	time.Sleep(100 * time.Millisecond)

	// 创建多个客户端
	numClients := 3
	clients := make([]*zmq.Socket, numClients)

	for i := 0; i < numClients; i++ {
		client, err := zmq.NewSocket(zmq.DEALER)
		if err != nil {
			t.Fatalf("Failed to create client %d: %v", i, err)
		}

		identity := fmt.Sprintf("client-%d-%d", i, time.Now().UnixNano())
		client.SetIdentity(identity)
		client.Connect(endpoint)
		clients[i] = client

		// 发送会话创建请求
		request := createTestSessionRequest()
		client.SendBytes(request, 0)

		// 接收响应
		response, _ := client.RecvBytes(0)
		sessionID, _ := parseSessionResponse(response)
		t.Logf("Client %d: Session %d created", i, sessionID)
	}

	// 验证所有 session 都创建了
	time.Sleep(100 * time.Millisecond)
	if sessionMgr.SessionCount() != int64(numClients) {
		t.Errorf("Expected %d sessions, got %d", numClients, sessionMgr.SessionCount())
	}

	// 推送事件，所有客户端都应收到
	blockEvent := createTestBlockEvent()
	sessionMgr.eventChan <- blockEvent

	time.Sleep(200 * time.Millisecond)

	// 每个客户端接收事件
	for i, client := range clients {
		client.SetRcvtimeo(1 * time.Second)
		data, err := client.RecvBytes(0)
		if err != nil {
			t.Errorf("Client %d failed to receive event: %v", i, err)
		} else {
			t.Logf("Client %d: Received %d bytes", i, len(data))
		}
	}

	// 关闭第一个客户端
	clients[0].Close()
	t.Log("Client 0 disconnected")

	// Wait for ZMQ to detect the disconnect (ROUTER needs time to notice)
	time.Sleep(100 * time.Millisecond)

	// 推送另一个事件触发断线检测
	anotherEvent := createTestBlockEvent()
	anotherEvent.Number = 888
	sessionMgr.eventChan <- anotherEvent

	// Wait for event processing and cleanup
	time.Sleep(500 * time.Millisecond)

	// 验证只剩 2 个 session
	if sessionMgr.SessionCount() != int64(numClients-1) {
		t.Errorf("Expected %d sessions after disconnect, got %d", numClients-1, sessionMgr.SessionCount())
	}

	// 清理其他客户端
	for i := 1; i < numClients; i++ {
		clients[i].Close()
	}
}

// TestZMQTracer_EventFilter 测试事件过滤
func TestZMQTracer_EventFilter(t *testing.T) {
	endpoint := fmt.Sprintf("ipc:///tmp/test_tracer_filter_%d.sock", time.Now().UnixNano())

	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{endpoint},
		QueueSize:     100,
		MaxSessions:   10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	if err := sessionMgr.Start(); err != nil {
		t.Fatalf("Failed to start session manager: %v", err)
	}
	defer sessionMgr.Stop()

	time.Sleep(100 * time.Millisecond)

	// 客户端连接
	client, _ := zmq.NewSocket(zmq.DEALER)
	client.SetIdentity("filter-test-client")
	client.Connect(endpoint)
	defer client.Close()

	// 发送带过滤器的会话请求（只接收 BlockEvent）
	// EventType values: TxStart=1, TxEnd=2, Enter=3, Exit=4, BlockStart=5, BlockEnd=6, Log=7
	// To filter for BlockStart (value=5), use filterMask = 1 << 5 = 32
	request := createTestSessionRequestWithFilter(1 << 5) // BlockStart = 5
	client.SendBytes(request, 0)

	response, _ := client.RecvBytes(0)
	sessionID, _ := parseSessionResponse(response)
	t.Logf("Session %d created with filter", sessionID)

	time.Sleep(100 * time.Millisecond)

	// 发送不同类型的事件
	txEvent := createTestTxStartEvent()
	sessionMgr.eventChan <- txEvent

	blockEvent := createTestBlockEvent()
	sessionMgr.eventChan <- blockEvent

	// 客户端应该只收到 BlockEvent
	client.SetRcvtimeo(1 * time.Second)
	data, err := client.RecvBytes(0)
	if err != nil {
		t.Fatal("Should receive block event")
	}

	event, _ := parseSBEEvent(data)
	if event.(*ethereum_tracing.BlockEvent) == nil {
		t.Error("Expected BlockEvent, got other type")
	}

	// 尝试接收第二个事件（TxStartEvent 应该被过滤）
	client.SetRcvtimeo(500 * time.Millisecond)
	_, err = client.RecvBytes(0)
	if err == nil {
		t.Error("Should not receive TxStartEvent (filtered out)")
	}
}

// Helper functions

func createTestSessionRequest() []byte {
	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0, // 接收所有事件
		ClientID0:   12345,
		ClientID1:   67890,
	}

	encoder := NewSBEEventEncoder()
	data, err := encoder.Encode(req)
	if err != nil {
		panic(err)
	}
	return data
}

func createTestSessionRequestWithFilter(filterMask uint32) []byte {
	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  filterMask,
		ClientID0:   12345,
		ClientID1:   67890,
	}

	encoder := NewSBEEventEncoder()
	data, err := encoder.Encode(req)
	if err != nil {
		panic(err)
	}
	return data
}

func parseSessionResponse(data []byte) (uint32, error) {
	marshaller := ethereum_tracing.NewSbeGoMarshaller()
	reader := bytes.NewReader(data)

	var header ethereum_tracing.SbeGoMessageHeader
	if err := header.Decode(marshaller, reader); err != nil {
		return 0, err
	}

	var resp ethereum_tracing.SessionCreateResponse
	if err := resp.Decode(marshaller, reader, header.Version, header.BlockLength, true); err != nil {
		return 0, err
	}

	return resp.SessionID, nil
}

func createTestBlockEvent() *ethereum_tracing.BlockEvent {
	var hash, parentHash [32]uint8
	for i := range hash {
		hash[i] = byte(i)
		parentHash[i] = byte(i + 1)
	}

	return &ethereum_tracing.BlockEvent{
		EventType:      ethereum_tracing.EventType.BlockStart,
		Timestamp:      uint64(time.Now().UnixNano()),
		Number:         12345,
		Hash:           hash,
		ParentHash:     parentHash,
		BlockTimestamp: uint64(time.Now().Unix()),
	}
}

func createTestTxStartEvent() *ethereum_tracing.TxStartEvent {
	var txHash [32]uint8
	var from, to [20]uint8
	for i := range txHash {
		txHash[i] = byte(i)
	}
	for i := range from {
		from[i] = byte(i)
		to[i] = byte(i + 1)
	}

	return &ethereum_tracing.TxStartEvent{
		EventType: ethereum_tracing.EventType.TxStart,
		Timestamp: uint64(time.Now().UnixNano()),
		TxHash:    txHash,
		From:      from,
		To:        to,
		GasLimit:  100000,
		Nonce:     42,
	}
}

func parseBlockEvent(data []byte) (*ethereum_tracing.BlockEvent, error) {
	marshaller := ethereum_tracing.NewSbeGoMarshaller()
	reader := bytes.NewReader(data)

	var header ethereum_tracing.SbeGoMessageHeader
	if err := header.Decode(marshaller, reader); err != nil {
		return nil, err
	}

	var event ethereum_tracing.BlockEvent
	if err := event.Decode(marshaller, reader, header.Version, header.BlockLength, true); err != nil {
		return nil, err
	}

	return &event, nil
}

// TestZMQTracer_MultipleSockets tests multiple ROUTER socket endpoints
// Flow:
// 1. Session manager listens on multiple endpoints
// 2. Clients connect to different endpoints
// 3. Each client creates session and receives events
// 4. Verify all clients receive events via their respective sockets
func TestZMQTracer_MultipleSockets(t *testing.T) {
	// Create multiple endpoints
	timestamp := time.Now().UnixNano()
	endpoint1 := fmt.Sprintf("ipc:///tmp/test_tracer_sock1_%d.sock", timestamp)
	endpoint2 := fmt.Sprintf("ipc:///tmp/test_tracer_sock2_%d.sock", timestamp)

	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{endpoint1, endpoint2},
		QueueSize:     100,
		MaxSessions:   10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	if err := sessionMgr.Start(); err != nil {
		t.Fatalf("Failed to start session manager: %v", err)
	}
	defer sessionMgr.Stop()

	// Wait for sockets to bind
	time.Sleep(100 * time.Millisecond)

	// Client 1 connects to endpoint 1 (socket index 0)
	client1, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		t.Fatalf("Failed to create client1: %v", err)
	}
	client1Identity := fmt.Sprintf("client1-ep1-%d", timestamp)
	client1.SetIdentity(client1Identity)
	client1.Connect(endpoint1)
	defer client1.Close()

	// Client 2 connects to endpoint 2 (socket index 1)
	client2, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		t.Fatalf("Failed to create client2: %v", err)
	}
	client2Identity := fmt.Sprintf("client2-ep2-%d", timestamp)
	client2.SetIdentity(client2Identity)
	client2.Connect(endpoint2)
	defer client2.Close()

	// Client 3 also connects to endpoint 1 (same socket as client1)
	client3, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		t.Fatalf("Failed to create client3: %v", err)
	}
	client3Identity := fmt.Sprintf("client3-ep1-%d", timestamp)
	client3.SetIdentity(client3Identity)
	client3.Connect(endpoint1)
	defer client3.Close()

	// Wait for connections
	time.Sleep(100 * time.Millisecond)

	// Create sessions
	request := createTestSessionRequest()

	// Client 1 creates session on socket 0
	client1.SendBytes(request, 0)
	response1, _ := client1.RecvBytes(0)
	sessionID1, _ := parseSessionResponse(response1)
	t.Logf("Client1: Session %d created on endpoint1 (socket 0)", sessionID1)

	// Client 2 creates session on socket 1
	client2.SendBytes(request, 0)
	response2, _ := client2.RecvBytes(0)
	sessionID2, _ := parseSessionResponse(response2)
	t.Logf("Client2: Session %d created on endpoint2 (socket 1)", sessionID2)

	// Client 3 creates session on socket 0 (same as client 1)
	client3.SendBytes(request, 0)
	response3, _ := client3.RecvBytes(0)
	sessionID3, _ := parseSessionResponse(response3)
	t.Logf("Client3: Session %d created on endpoint1 (socket 0)", sessionID3)

	// Wait for session creation
	time.Sleep(100 * time.Millisecond)

	// Verify all sessions created
	if sessionMgr.SessionCount() != 3 {
		t.Errorf("Expected 3 sessions, got %d", sessionMgr.SessionCount())
	}

	// Verify sessions are on correct sockets
	sessionMgr.mu.RLock()
	for _, session := range sessionMgr.sessions {
		switch session.ID {
		case sessionID1, sessionID3:
			if session.socketIndex != 0 {
				t.Errorf("Session %d should be on socket 0, got %d", session.ID, session.socketIndex)
			}
		case sessionID2:
			if session.socketIndex != 1 {
				t.Errorf("Session %d should be on socket 1, got %d", session.ID, session.socketIndex)
			}
		}
	}
	sessionMgr.mu.RUnlock()

	// Send block event to all sessions
	blockEvent := createTestBlockEvent()
	blockEvent.Number = 12345
	sessionMgr.eventChan <- blockEvent

	// Wait for event distribution
	time.Sleep(200 * time.Millisecond)

	// All 3 clients should receive the event
	for i, client := range []*zmq.Socket{client1, client2, client3} {
		client.SetRcvtimeo(1 * time.Second)
		data, err := client.RecvBytes(0)
		if err != nil {
			t.Errorf("Client %d failed to receive event: %v", i+1, err)
			continue
		}

		event, err := parseBlockEvent(data)
		if err != nil {
			t.Errorf("Client %d failed to parse event: %v", i+1, err)
			continue
		}

		if event.Number != 12345 {
			t.Errorf("Client %d: Expected block number 12345, got %d", i+1, event.Number)
		} else {
			t.Logf("Client %d: Received block event successfully (block=%d)", i+1, event.Number)
		}
	}

	// Disconnect client1 and verify session cleanup
	client1.Close()
	t.Log("Client1 disconnected")

	// Wait for disconnect detection
	time.Sleep(100 * time.Millisecond)

	// Send another event to trigger disconnect detection
	anotherEvent := createTestBlockEvent()
	anotherEvent.Number = 99999
	sessionMgr.eventChan <- anotherEvent

	// Wait for cleanup
	time.Sleep(300 * time.Millisecond)

	// Should have 2 sessions remaining (client2 and client3)
	if sessionMgr.SessionCount() != 2 {
		t.Errorf("Expected 2 sessions after client1 disconnect, got %d", sessionMgr.SessionCount())
	}

	// Client2 and Client3 should still receive events
	client2.SetRcvtimeo(1 * time.Second)
	data2, err := client2.RecvBytes(0)
	if err != nil {
		t.Errorf("Client2 failed to receive event after client1 disconnect: %v", err)
	} else {
		event2, _ := parseBlockEvent(data2)
		t.Logf("Client2: Still receiving events (block=%d)", event2.Number)
	}

	client3.SetRcvtimeo(1 * time.Second)
	data3, err := client3.RecvBytes(0)
	if err != nil {
		t.Errorf("Client3 failed to receive event after client1 disconnect: %v", err)
	} else {
		event3, _ := parseBlockEvent(data3)
		t.Logf("Client3: Still receiving events (block=%d)", event3.Number)
	}
}

func parseSBEEvent(data []byte) (interface{}, error) {
	marshaller := ethereum_tracing.NewSbeGoMarshaller()
	reader := bytes.NewReader(data)

	var header ethereum_tracing.SbeGoMessageHeader
	if err := header.Decode(marshaller, reader); err != nil {
		return nil, err
	}

	// 根据模板ID解析不同类型
	switch header.TemplateId {
	case new(ethereum_tracing.BlockEvent).SbeTemplateId():
		var event ethereum_tracing.BlockEvent
		err := event.Decode(marshaller, reader, header.Version, header.BlockLength, true)
		return &event, err
	case new(ethereum_tracing.TxStartEvent).SbeTemplateId():
		var event ethereum_tracing.TxStartEvent
		err := event.Decode(marshaller, reader, header.Version, header.BlockLength, true)
		return &event, err
	default:
		return nil, fmt.Errorf("unknown template id: %d", header.TemplateId)
	}
}
