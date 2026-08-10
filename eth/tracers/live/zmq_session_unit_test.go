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
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/eth/tracers/live/sbe_generated/ethereum_tracing"
)

// TestZMQSessionManager_CreateSession 测试会话创建
func TestZMQSessionManager_CreateSession(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:     100,
		MaxSessions:   10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// 模拟创建会话
	identity1 := []byte("client-1")
	req1 := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0,
	}

	session1, err := sessionMgr.createSession(identity1, req1, 0)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session1.ID != 0 {
		t.Errorf("Expected session ID 0, got %d", session1.ID)
	}

	if sessionMgr.SessionCount() != 1 {
		t.Errorf("Expected 1 session, got %d", sessionMgr.SessionCount())
	}

	// 创建第二个会话
	identity2 := []byte("client-2")
	session2, err := sessionMgr.createSession(identity2, req1, 0)
	if err != nil {
		t.Fatalf("Failed to create second session: %v", err)
	}

	if session2.ID != 1 {
		t.Errorf("Expected session ID 1, got %d", session2.ID)
	}

	if sessionMgr.SessionCount() != 2 {
		t.Errorf("Expected 2 sessions, got %d", sessionMgr.SessionCount())
	}

	// 重复创建相同 identity 应该返回已有会话
	session1Again, err := sessionMgr.createSession(identity1, req1, 0)
	if err != nil {
		t.Fatalf("Failed to recreate session: %v", err)
	}

	if session1Again.ID != session1.ID {
		t.Errorf("Expected same session ID %d, got %d", session1.ID, session1Again.ID)
	}

	if sessionMgr.SessionCount() != 2 {
		t.Errorf("Expected 2 sessions after recreate, got %d", sessionMgr.SessionCount())
	}
}

// TestZMQSessionManager_CleanupSession 测试会话清理
func TestZMQSessionManager_CleanupSession(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// 创建多个会话
	identity1 := []byte("client-1")
	identity2 := []byte("client-2")
	identity3 := []byte("client-3")

	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0,
	}

	sessionMgr.createSession(identity1, req, 0)
	sessionMgr.createSession(identity2, req, 0)
	sessionMgr.createSession(identity3, req, 0)

	if sessionMgr.SessionCount() != 3 {
		t.Errorf("Expected 3 sessions, got %d", sessionMgr.SessionCount())
	}

	// 清理一个会话
	sessionMgr.cleanupSession(identity2)

	if sessionMgr.SessionCount() != 2 {
		t.Errorf("Expected 2 sessions after cleanup, got %d", sessionMgr.SessionCount())
	}

	// 还有 2 个会话，IsEmpty 应该返回 false
	if sessionMgr.IsEmpty() {
		t.Errorf("Expected IsEmpty to return false (have %d sessions), but got true", sessionMgr.SessionCount())
	}

	// 清理不存在的 identity 应该无影响
	sessionMgr.cleanupSession([]byte("non-existent"))

	if sessionMgr.SessionCount() != 2 {
		t.Errorf("Expected 2 sessions after non-existent cleanup, got %d", sessionMgr.SessionCount())
	}

	// 清理剩余会话
	sessionMgr.cleanupSession(identity1)
	sessionMgr.cleanupSession(identity3)

	if sessionMgr.SessionCount() != 0 {
		t.Errorf("Expected 0 sessions, got %d", sessionMgr.SessionCount())
	}

	if !sessionMgr.IsEmpty() {
		t.Error("Expected IsEmpty to return true")
	}
}

// TestZMQSessionManager_MaxSessions 测试最大会话数限制
func TestZMQSessionManager_MaxSessions(t *testing.T) {
	maxSessions := 3
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  maxSessions,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0,
	}

	// 创建达到上限的会话
	for i := 0; i < maxSessions; i++ {
		identity := []byte{byte(i)}
		_, err := sessionMgr.createSession(identity, req, 0)
		if err != nil {
			t.Fatalf("Failed to create session %d: %v", i, err)
		}
	}

	// 尝试创建超出上限的会话
	identity := []byte{byte(maxSessions)}
	_, err = sessionMgr.createSession(identity, req, 0)
	if err == nil {
		t.Error("Expected error when creating session beyond max limit")
	}
	if err != ErrMaxSessionsReached {
		t.Errorf("Expected ErrMaxSessionsReached, got %v", err)
	}
}

// TestZMQSessionManager_EventFilter 测试事件过滤
func TestZMQSessionManager_EventFilter(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// 创建会话，只接收 BlockEvent
	// EventType: TxStart=1, TxEnd=2, Enter=3, Exit=4, BlockStart=5, BlockEnd=6, Log=7
	// 要过滤 BlockStart (值=5)，filterMask 应该是 1 << 5 = 32
	identity := []byte("client-1")
	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  1 << 5, // BlockStart = 5
	}

	session, err := sessionMgr.createSession(identity, req, 0)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// 测试事件过滤
	txEvent := &ethereum_tracing.TxStartEvent{
		EventType: ethereum_tracing.EventType.TxStart,
	}

	blockEvent := &ethereum_tracing.BlockEvent{
		EventType: ethereum_tracing.EventType.BlockStart,
	}

	// TxStart 应该被过滤掉
	if session.shouldProcessEvent(txEvent) {
		t.Error("TxStart should be filtered out")
	}

	// BlockStart 应该通过
	if !session.shouldProcessEvent(blockEvent) {
		t.Error("BlockStart should pass filter")
	}
}

// TestZMQSessionManager_ConcurrentAccess 测试并发访问
func TestZMQSessionManager_ConcurrentAccess(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  100,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	numSessions := 50
	var wg sync.WaitGroup

	// 并发创建会话
	for i := 0; i < numSessions; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			identity := []byte{byte(id % 256), byte(id / 256)}
			req := &ethereum_tracing.SessionCreateRequest{
				MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
				FilterMask:  0,
			}
			_, err := sessionMgr.createSession(identity, req, 0)
			if err != nil {
				t.Errorf("Failed to create session %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	if sessionMgr.SessionCount() != int64(numSessions) {
		t.Errorf("Expected %d sessions, got %d", numSessions, sessionMgr.SessionCount())
	}

	// 并发清理会话
	for i := 0; i < numSessions; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			identity := []byte{byte(id % 256), byte(id / 256)}
			sessionMgr.cleanupSession(identity)
		}(i)
	}

	wg.Wait()

	if sessionMgr.SessionCount() != 0 {
		t.Errorf("Expected 0 sessions after cleanup, got %d", sessionMgr.SessionCount())
	}
}

// TestZMQSessionManager_EventBroadcast 测试事件广播（不使用 ZeroMQ）
func TestZMQSessionManager_EventBroadcast(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// 创建会话
	identity := []byte("test-client")
	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0, // 接收所有事件
	}

	session, err := sessionMgr.createSession(identity, req, 0)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// 测试 shouldProcessEvent 逻辑
	testCases := []struct {
		name      string
		event     SBEEvent
		shouldPro bool
	}{
		{
			name: "TxStartEvent",
			event: &ethereum_tracing.TxStartEvent{
				EventType: ethereum_tracing.EventType.TxStart,
			},
			shouldPro: true,
		},
		{
			name: "BlockEvent",
			event: &ethereum_tracing.BlockEvent{
				EventType: ethereum_tracing.EventType.BlockStart,
			},
			shouldPro: true,
		},
		{
			name: "TxEndEvent",
			event: &ethereum_tracing.TxEndEvent{
				EventType: ethereum_tracing.EventType.TxEnd,
			},
			shouldPro: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := session.shouldProcessEvent(tc.event)
			if result != tc.shouldPro {
				t.Errorf("Expected shouldProcessEvent to return %v, got %v", tc.shouldPro, result)
			}
		})
	}
}

// TestZMQSessionManager_IsEmpty 测试 IsEmpty 快速检查
func TestZMQSessionManager_IsEmpty(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// 初始应该是空的
	if !sessionMgr.IsEmpty() {
		t.Error("Expected IsEmpty to return true initially")
	}

	// 创建会话后不应该是空的
	identity := []byte("client-1")
	req := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0,
	}

	sessionMgr.createSession(identity, req, 0)

	if sessionMgr.IsEmpty() {
		t.Error("Expected IsEmpty to return false after creating session")
	}

	// 清理后应该是空的
	sessionMgr.cleanupSession(identity)

	if !sessionMgr.IsEmpty() {
		t.Error("Expected IsEmpty to return true after cleanup")
	}
}

// TestSession_ShouldProcessEvent 测试事件过滤逻辑
func TestSession_ShouldProcessEvent(t *testing.T) {
	// EventType values: TxStart=1, TxEnd=2, Enter=3, Exit=4, BlockStart=5, BlockEnd=6, Log=7
	// To filter event with value V, use filterMask = 1 << V
	testCases := []struct {
		name       string
		filterMask uint32
		eventType  ethereum_tracing.EventTypeEnum
		expected   bool
	}{
		{
			name:       "No filter - accept all",
			filterMask: 0,
			eventType:  ethereum_tracing.EventType.TxStart,
			expected:   true,
		},
		{
			name:       "Filter TxStart only",
			filterMask: 1 << 1, // TxStart = 1
			eventType:  ethereum_tracing.EventType.TxStart,
			expected:   true,
		},
		{
			name:       "Filter TxStart - BlockStart blocked",
			filterMask: 1 << 1, // TxStart = 1
			eventType:  ethereum_tracing.EventType.BlockStart,
			expected:   false,
		},
		{
			name:       "Multiple filters",
			filterMask: (1 << 1) | (1 << 5), // TxStart and BlockStart
			eventType:  ethereum_tracing.EventType.TxStart,
			expected:   true,
		},
		{
			name:       "Multiple filters - TxEnd blocked",
			filterMask: (1 << 1) | (1 << 5), // TxStart and BlockStart
			eventType:  ethereum_tracing.EventType.TxEnd,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			session := &Session{
				ID:         0,
				identity:   []byte("test"),
				filterMask: tc.filterMask,
			}

			var event SBEEvent
			switch tc.eventType {
			case ethereum_tracing.EventType.TxStart:
				event = &ethereum_tracing.TxStartEvent{EventType: tc.eventType}
			case ethereum_tracing.EventType.TxEnd:
				event = &ethereum_tracing.TxEndEvent{EventType: tc.eventType}
			case ethereum_tracing.EventType.BlockStart:
				event = &ethereum_tracing.BlockEvent{EventType: tc.eventType}
			}

			result := session.shouldProcessEvent(event)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
// TestZMQSessionManager_PerEventTypeCount tests per-event-type session counting
func TestZMQSessionManager_PerEventTypeCount(t *testing.T) {
	sessionMgr, err := NewZMQSessionManager(&ZMQSessionManagerConfig{
		BindEndpoints: []string{"ipc:///tmp/test.sock"},
		QueueSize:    100,
		MaxSessions:  10,
	})
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	// Initially, all event type counts should be 0
	for eventType := sessionMgr.eventTypeStart; eventType < sessionMgr.eventTypeEnd; eventType++ {
		if sessionMgr.HasSessionsForEventType(ethereum_tracing.EventTypeEnum(eventType)) {
			t.Errorf("Expected no sessions for event type %d", eventType)
		}
	}

	// Create session subscribing to TxStart only (eventType=1)
	identity1 := []byte("client-1")
	req1 := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  1 << 1, // TxStart = 1
	}
	_, err = sessionMgr.createSession(identity1, req1, 0)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Verify counts
	if !sessionMgr.HasSessionsForEventType(ethereum_tracing.EventType.TxStart) {
		t.Error("Expected sessions for TxStart")
	}
	if sessionMgr.HasSessionsForEventType(ethereum_tracing.EventType.TxEnd) {
		t.Error("Expected no sessions for TxEnd")
	}

	// Create session subscribing to all events (filterMask=0)
	identity2 := []byte("client-2")
	req2 := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  0, // All events
	}
	_, err = sessionMgr.createSession(identity2, req2, 0)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Now all event types should have at least one session
	for eventType := sessionMgr.eventTypeStart; eventType < sessionMgr.eventTypeEnd; eventType++ {
		if !sessionMgr.HasSessionsForEventType(ethereum_tracing.EventTypeEnum(eventType)) {
			t.Errorf("Expected sessions for event type %d (client-2 subscribes all)", eventType)
		}
	}

	// TxStart should have 2 sessions (client-1 + client-2)
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load() != 2 {
		t.Errorf("Expected 2 sessions for TxStart, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load())
	}

	// TxEnd should have 1 session (client-2 only)
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxEnd].Load() != 1 {
		t.Errorf("Expected 1 session for TxEnd, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxEnd].Load())
	}

	// Total session count should be 2
	if sessionMgr.SessionCount() != 2 {
		t.Errorf("Expected total 2 sessions, got %d", sessionMgr.SessionCount())
	}

	// Remove client-1 (subscribed to TxStart only)
	sessionMgr.cleanupSession(identity1)

	// TxStart should now have 1 session (client-2)
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load() != 1 {
		t.Errorf("Expected 1 session for TxStart after cleanup, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load())
	}

	// TxEnd should still have 1 session (client-2)
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxEnd].Load() != 1 {
		t.Errorf("Expected 1 session for TxEnd, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxEnd].Load())
	}

	// Total should be 1
	if sessionMgr.SessionCount() != 1 {
		t.Errorf("Expected total 1 session, got %d", sessionMgr.SessionCount())
	}

	// Remove client-2 (subscribed to all)
	sessionMgr.cleanupSession(identity2)

	// All counts should be 0
	for eventType := 0; eventType < sessionMgr.eventTypeEnd; eventType++ {
		if sessionMgr.sessionCounts[eventType].Load() != 0 {
			t.Errorf("Expected 0 sessions for index %d, got %d", eventType, sessionMgr.sessionCounts[eventType].Load())
		}
	}

	// Test filter mask update
	identity3 := []byte("client-3")
	req3 := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  1 << 1, // TxStart only
	}
	session3, err := sessionMgr.createSession(identity3, req3, 0)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Verify TxStart has 1 session
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load() != 1 {
		t.Errorf("Expected 1 session for TxStart, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load())
	}

	// Update filter mask to TxEnd (eventType=2)
	req3Update := &ethereum_tracing.SessionCreateRequest{
		MessageType: ethereum_tracing.MessageType.SessionCreateRequest,
		FilterMask:  1 << 2, // TxEnd only
	}
	_, err = sessionMgr.createSession(identity3, req3Update, 0) // Same identity updates existing session
	if err != nil {
		t.Fatalf("Failed to update session: %v", err)
	}

	// TxStart should now have 0 sessions
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load() != 0 {
		t.Errorf("Expected 0 sessions for TxStart after update, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxStart].Load())
	}

	// TxEnd should have 1 session
	if sessionMgr.sessionCounts[ethereum_tracing.EventType.TxEnd].Load() != 1 {
		t.Errorf("Expected 1 session for TxEnd after update, got %d", sessionMgr.sessionCounts[ethereum_tracing.EventType.TxEnd].Load())
	}

	// Total should still be 1 (session count unchanged)
	if sessionMgr.SessionCount() != 1 {
		t.Errorf("Expected total 1 session, got %d", sessionMgr.SessionCount())
	}

	t.Logf("Session ID %d filter updated successfully", session3.ID)
}
