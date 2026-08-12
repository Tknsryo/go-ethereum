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
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/eth/tracers/live/sbe_generated/ethereum_tracing"
	"github.com/ethereum/go-ethereum/log"
	zmq "gopkg.in/pebbe/zmq4.v1"
)

var (
	ErrSessionManagerClosed = errors.New("session manager is closed")
	ErrMaxSessionsReached   = errors.New("maximum number of sessions reached")
)

// getEventTypeRange extracts event type boundaries from SBE-generated EventType
// Returns (min, max+1) where max+1 is the exclusive end
func getEventTypeRange() (int, int) {
	val := reflect.ValueOf(ethereum_tracing.EventType)
	typ := reflect.TypeOf(ethereum_tracing.EventType)

	min := int(^uint8(0)) // Max uint8
	max := 0

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		// Skip NullValue field
		if field.Name == "NullValue" {
			continue
		}

		value := int(val.Field(i).Uint())
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return min, max + 1 // Return exclusive end
}

// ZMQSessionManager manages client sessions and distributes events
type ZMQSessionManager struct {
	routerSockets     []*zmq.Socket // Multiple ROUTER sockets for all communication
	bindEndpoints     []string      // ROUTER socket endpoints
	sessions          map[uint32]*Session
	identityToSession map[string]*Session // Identity -> Session mapping
	eventChan         chan SBEEvent       // Ingress channel from hooks
	shutdownChan      chan struct{}       // Global shutdown signal
	nextSessionID     uint32
	maxSessions       int
	queueSize         int
	sendTimeout       time.Duration // Send timeout for POLLOUT check
	wg                sync.WaitGroup
	mu                sync.RWMutex
	closed            bool
	sessionCounts     []atomic.Int64   // Per-event-type session count (dynamic size)
	eventTypeStart    int              // First valid event type (from reflection)
	eventTypeEnd      int              // One past last event type (from reflection)
	encoder           *SBEEventEncoder // Shared encoder
}

// Session represents a single client session
type Session struct {
	ID          uint32
	identity    []byte // Client identity (for ROUTER routing)
	filterMask  uint32 // Event type filter
	socketIndex int    // Index of ROUTER socket this session belongs to
	manager     *ZMQSessionManager

	// LogEvent specific filters (eth_getLogs style)
	// Using byte arrays for efficient comparison without hex conversion
	logAddressFilter map[[20]byte]bool   // nil/empty = accept all addresses
	logTopicsFilter  []map[[32]byte]bool // nil/empty = accept all topics, supports empty map as wildcard
}

// ZMQSessionManagerConfig holds configuration for session manager
type ZMQSessionManagerConfig struct {
	BindEndpoints []string // ZMQ ROUTER endpoints (e.g., ["ipc:///tmp/geth.tracer.sock", "tcp://*:5555"])
	QueueSize     int      // Event queue size
	MaxSessions   int      // Maximum concurrent sessions (0 = unlimited)
	SendTimeout   int      // Send timeout in milliseconds (0 = default 2000ms)
	// Note: Lower SendTimeout detects disconnected clients faster but may timeout
	// on slow networks. For local IPC, 100-500ms is fine. For remote TCP clients
	// over unreliable networks, use 2000-5000ms or higher.
}

// NewZMQSessionManager creates a new session manager
func NewZMQSessionManager(cfg *ZMQSessionManagerConfig) (*ZMQSessionManager, error) {
	// Set defaults
	if cfg.QueueSize == 0 {
		cfg.QueueSize = 1000
	}
	if cfg.MaxSessions == 0 {
		cfg.MaxSessions = 100
	}

	// Set send timeout (default 2 seconds for network-friendly operation)
	sendTimeout := time.Duration(cfg.SendTimeout) * time.Millisecond
	if cfg.SendTimeout == 0 {
		sendTimeout = 2 * time.Second
	}

	// Get event type range from SBE-generated EventType
	eventTypeStart, eventTypeEnd := getEventTypeRange()

	// Create session counts slice: [0] = total, [eventTypeStart..eventTypeEnd) = per-event-type
	sessionCounts := make([]atomic.Int64, eventTypeEnd)

	// Create ROUTER sockets for each endpoint
	routerSockets := make([]*zmq.Socket, 0, len(cfg.BindEndpoints))
	bindEndpoints := make([]string, 0, len(cfg.BindEndpoints))

	for _, endpoint := range cfg.BindEndpoints {
		socket, err := zmq.NewSocket(zmq.ROUTER)
		if err != nil {
			// Clean up already created sockets
			for _, s := range routerSockets {
				s.Close()
			}
			return nil, fmt.Errorf("failed to create ROUTER socket for %s: %w", endpoint, err)
		}

		// Set ROUTER_MANDATORY to detect disconnected clients
		if err := socket.SetRouterMandatory(1); err != nil {
			socket.Close()
			for _, s := range routerSockets {
				s.Close()
			}
			return nil, fmt.Errorf("failed to set ROUTER_MANDATORY for %s: %w", endpoint, err)
		}

		routerSockets = append(routerSockets, socket)
		bindEndpoints = append(bindEndpoints, endpoint)
	}

	return &ZMQSessionManager{
		routerSockets:     routerSockets,
		bindEndpoints:     bindEndpoints,
		sessions:          make(map[uint32]*Session),
		identityToSession: make(map[string]*Session),
		eventChan:         make(chan SBEEvent, cfg.QueueSize),
		shutdownChan:      make(chan struct{}),
		maxSessions:       cfg.MaxSessions,
		queueSize:         cfg.QueueSize,
		sendTimeout:       sendTimeout,
		sessionCounts:     sessionCounts,
		eventTypeStart:    eventTypeStart,
		eventTypeEnd:      eventTypeEnd,
		encoder:           NewSBEEventEncoder(),
	}, nil
}

// Start starts the session manager
func (m *ZMQSessionManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrSessionManagerClosed
	}

	// Bind all ROUTER sockets
	for i, socket := range m.routerSockets {
		endpoint := m.bindEndpoints[i]

		// Bind socket
		if err := socket.Bind(endpoint); err != nil {
			return fmt.Errorf("failed to bind ROUTER socket to %s: %w", endpoint, err)
		}

		log.Info("ZMQ ROUTER socket bound", "endpoint", endpoint, "index", i)
	}

	// Start message handler (handles session creation requests)
	m.wg.Add(1)
	go m.handleMessages()

	// Start event distribution loop
	m.wg.Add(1)
	go m.distributeEvents()

	log.Info("ZMQ tracer session manager started", "endpoints", m.bindEndpoints)
	return nil
}

// handleMessages handles incoming messages (session creation and client requests)
// Uses zmq.Poller to monitor multiple ROUTER sockets
func (m *ZMQSessionManager) handleMessages() {
	defer m.wg.Done()

	// Create poller and add all ROUTER sockets
	poller := zmq.NewPoller()
	for _, socket := range m.routerSockets {
		poller.Add(socket, zmq.POLLIN)
	}

	for {
		select {
		case <-m.shutdownChan:
			return
		default:
			// Poll all sockets with timeout for responsive shutdown
			sockets, err := poller.Poll(100 * time.Millisecond)
			if err != nil {
				if err.Error() == "resource temporarily unavailable" {
					continue // Timeout, check shutdown
				}
				log.Error("Poll failed", "error", err)
				continue
			}

			// Process each socket that has data
			for _, polled := range sockets {
				// Find socket index
				socketIdx := m.findSocketIndex(polled.Socket)
				if socketIdx < 0 {
					log.Error("Unknown socket in poll result")
					continue
				}

				// Handle message from this socket
				m.handleSocketMessage(socketIdx)
			}
		}
	}
}

// findSocketIndex finds the index of a socket in routerSockets
func (m *ZMQSessionManager) findSocketIndex(socket *zmq.Socket) int {
	for i, s := range m.routerSockets {
		if s == socket {
			return i
		}
	}
	return -1
}

// handleSocketMessage handles a message from a specific ROUTER socket
func (m *ZMQSessionManager) handleSocketMessage(socketIdx int) {
	socket := m.routerSockets[socketIdx]

	// ROUTER receives from DEALER: [identity][data]
	// Note: DEALER doesn't add delimiter for single-frame messages
	identity, err := socket.RecvBytes(0)
	if err != nil {
		if err.Error() == "socket is closed" {
			return
		}
		log.Error("Failed to receive identity", "socketIndex", socketIdx, "error", err)
		return
	}

	// Check if there's more data (should be true for DEALER messages)
	more, err := socket.GetRcvmore()
	if err != nil {
		log.Error("Failed to check Rcvmore", "socketIndex", socketIdx, "error", err)
		return
	}

	if !more {
		log.Warn("Received message without data part", "socketIndex", socketIdx)
		return
	}

	// Receive data
	data, err := socket.RecvBytes(0)
	if err != nil {
		log.Error("Failed to receive data", "socketIndex", socketIdx, "error", err)
		return
	}

	// Process message (pass socket index for session creation)
	response := m.processMessage(identity, data, socketIdx)

	// Send response: [identity][data]
	if len(response) > 0 {
		if err := m.sendToClient(identity, response, socketIdx); err != nil {
			log.Error("Failed to send response", "socketIndex", socketIdx, "error", err)
		}
	}
}

// processMessage processes an incoming message
func (m *ZMQSessionManager) processMessage(identity []byte, data []byte, socketIdx int) []byte {
	// Decode message
	marshaller := ethereum_tracing.NewSbeGoMarshaller()
	reader := bytes.NewReader(data)

	// Decode message header
	var header ethereum_tracing.SbeGoMessageHeader
	if err := header.Decode(marshaller, reader); err != nil {
		return m.encodeErrorResponse("failed to decode header: " + err.Error())
	}

	// Check message type from templateId
	switch header.TemplateId {
	case new(ethereum_tracing.SessionCreateRequest).SbeTemplateId():
		return m.handleSessionCreate(identity, data, socketIdx)
	default:
		return m.encodeErrorResponse("unknown message type")
	}
}

// handleSessionCreate handles session creation request
func (m *ZMQSessionManager) handleSessionCreate(identity []byte, data []byte, socketIdx int) []byte {
	marshaller := ethereum_tracing.NewSbeGoMarshaller()
	reader := bytes.NewReader(data)

	// Decode header
	var header ethereum_tracing.SbeGoMessageHeader
	if err := header.Decode(marshaller, reader); err != nil {
		return m.encodeErrorResponse("failed to decode header: " + err.Error())
	}

	// Decode SessionCreateRequest
	// Note: Range check disabled because Hash/Address are raw byte arrays
	// where each byte can be 0-255 (no null value reservation needed)
	var req ethereum_tracing.SessionCreateRequest
	if err := req.Decode(marshaller, reader, header.Version, header.BlockLength, false); err != nil {
		return m.encodeErrorResponse("failed to decode request: " + err.Error())
	}

	// Validate request
	if req.MessageType != ethereum_tracing.MessageType.SessionCreateRequest {
		return m.encodeErrorResponse("invalid message type")
	}

	// Create session
	session, err := m.createSession(identity, &req, socketIdx)
	if err != nil {
		return m.encodeErrorResponse(err.Error())
	}

	// Encode success response (return endpoint for this socket)
	return m.encodeSuccessResponse(session.ID, m.bindEndpoints[socketIdx])
}

// incrementSessionCounts increments session counts for event types in filter mask
// filterMask=0 means subscribe to all events
func (m *ZMQSessionManager) incrementSessionCounts(filterMask uint32) {
	// Always increment total count
	m.sessionCounts[0].Add(1)

	if filterMask == 0 {
		// Subscribe to all event types
		for eventType := m.eventTypeStart; eventType < m.eventTypeEnd; eventType++ {
			m.sessionCounts[eventType].Add(1)
		}
	} else {
		// Subscribe to specific event types
		for eventType := m.eventTypeStart; eventType < m.eventTypeEnd; eventType++ {
			if (filterMask & (1 << eventType)) != 0 {
				m.sessionCounts[eventType].Add(1)
			}
		}
	}
}

// decrementSessionCounts decrements session counts for event types in filter mask
func (m *ZMQSessionManager) decrementSessionCounts(filterMask uint32) {
	// Always decrement total count
	m.sessionCounts[0].Add(-1)

	if filterMask == 0 {
		// Unsubscribe from all event types
		for eventType := m.eventTypeStart; eventType < m.eventTypeEnd; eventType++ {
			m.sessionCounts[eventType].Add(-1)
		}
	} else {
		// Unsubscribe from specific event types
		for eventType := m.eventTypeStart; eventType < m.eventTypeEnd; eventType++ {
			if (filterMask & (1 << eventType)) != 0 {
				m.sessionCounts[eventType].Add(-1)
			}
		}
	}
}

// updateSessionCounts adjusts session counts when filter mask changes
func (m *ZMQSessionManager) updateSessionCounts(oldMask, newMask uint32) {
	// Decrement old mask subscriptions
	m.decrementSessionCounts(oldMask)
	// Increment new mask subscriptions
	m.incrementSessionCounts(newMask)
}

// createSession creates a new session
func (m *ZMQSessionManager) createSession(identity []byte, req *ethereum_tracing.SessionCreateRequest, socketIdx int) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if closed
	if m.closed {
		return nil, ErrSessionManagerClosed
	}

	// Check max sessions
	if m.maxSessions > 0 && len(m.sessions) >= m.maxSessions {
		return nil, ErrMaxSessionsReached
	}

	// Check if identity already has a session
	identityKey := string(identity)
	if existingSession, exists := m.identityToSession[identityKey]; exists {
		// Update existing session's filter mask - need to adjust counts
		oldFilterMask := existingSession.filterMask
		existingSession.filterMask = req.FilterMask

		// Adjust event type counts
		m.updateSessionCounts(oldFilterMask, req.FilterMask)

		log.Info("Session updated",
			"sessionID", existingSession.ID,
			"filterMask", req.FilterMask,
			"socketIndex", existingSession.socketIndex)
		return existingSession, nil
	}

	// Generate session ID
	sessionID := m.nextSessionID
	m.nextSessionID++

	// Build log address filter (eth_getLogs style)
	var logAddressFilter map[[20]byte]bool
	if len(req.AddressFilter) > 0 {
		logAddressFilter = make(map[[20]byte]bool)
		for _, addr := range req.AddressFilter {
			logAddressFilter[addr.Address] = true
		}
	}

	// Build log topics filter (eth_getLogs style)
	// Each position in the outer slice represents a topic position
	// Empty map at position i means "match any topic at position i"
	var logTopicsFilter []map[[32]byte]bool
	if len(req.TopicsFilter) > 0 {
		logTopicsFilter = make([]map[[32]byte]bool, len(req.TopicsFilter))
		for i, topicFilter := range req.TopicsFilter {
			if len(topicFilter.Topics) > 0 {
				logTopicsFilter[i] = make(map[[32]byte]bool)
				for _, topic := range topicFilter.Topics {
					logTopicsFilter[i][topic.Topic] = true
				}
			}
			// Empty TopicsFilter[i].Topics means empty map = match any topic at position i
		}
	}

	// Create session
	session := &Session{
		ID:               sessionID,
		identity:         identity,
		filterMask:       req.FilterMask,
		socketIndex:      socketIdx,
		manager:          m,
		logAddressFilter: logAddressFilter,
		logTopicsFilter:  logTopicsFilter,
	}

	m.sessions[sessionID] = session
	m.identityToSession[identityKey] = session

	// Increment session counts based on filter mask
	m.incrementSessionCounts(req.FilterMask)

	log.Info("Session created",
		"sessionID", sessionID,
		"identity", fmt.Sprintf("%x", identity),
		"filterMask", req.FilterMask,
		"activeSessions", m.sessionCounts[0].Load())

	return session, nil
}

// distributeEvents distributes events to all active sessions
func (m *ZMQSessionManager) distributeEvents() {
	defer m.wg.Done()

	for {
		select {
		case event := <-m.eventChan:
			// Broadcast to all sessions
			m.broadcastEvent(event)

		case <-m.shutdownChan:
			return
		}
	}
}

// broadcastEvent sends event to all matching sessions
func (m *ZMQSessionManager) broadcastEvent(event SBEEvent) {
	// Encode event once (outside lock)
	encoded, err := m.encoder.Encode(event)
	if err != nil {
		log.Error("Failed to encode event", "error", err)
		return
	}

	// Get list of sessions to send to (with read lock)
	type sessionInfo struct {
		identity    []byte
		socketIndex int
		shouldProc  bool
	}

	m.mu.RLock()
	sessions := make([]sessionInfo, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, sessionInfo{
			identity:    session.identity,
			socketIndex: session.socketIndex,
			shouldProc:  session.shouldProcessEvent(event),
		})
	}
	m.mu.RUnlock()

	// Send to each session that matches filter (outside lock)
	for _, info := range sessions {
		if !info.shouldProc {
			continue
		}

		// Send via ROUTER: [identity][encoded_event]
		if err := m.sendToClient(info.identity, encoded, info.socketIndex); err != nil {
			// Client disconnected - cleanup session (async to avoid deadlock)
			go m.cleanupSession(info.identity)
		}
	}
}

// sendToClient sends a message to a specific client via ROUTER socket
// Returns error if client is disconnected (EHOSTUNREACH) or socket not writable
func (m *ZMQSessionManager) sendToClient(identity []byte, data []byte, socketIdx int) error {
	// Validate socket index
	if socketIdx < 0 || socketIdx >= len(m.routerSockets) {
		return fmt.Errorf("invalid socket index: %d", socketIdx)
	}

	socket := m.routerSockets[socketIdx]

	// Check if socket is writable using Poll with POLLOUT
	// This prevents blocking when client buffer is full or disconnected
	// Timeout is configurable (default 2s) for network-friendly operation
	poller := zmq.NewPoller()
	poller.Add(socket, zmq.POLLOUT)

	sockets, err := poller.Poll(m.sendTimeout)
	if err != nil {
		return fmt.Errorf("poll failed: %w", err)
	}

	// Check if socket is ready for writing
	if len(sockets) == 0 {
		// Socket not writable (timeout) - client likely disconnected or buffer full
		return fmt.Errorf("socket not writable (timeout after %v)", m.sendTimeout)
	}

	// Socket is ready, send the message
	// ROUTER send format: [identity][data]
	// Note: No delimiter needed for DEALER clients
	// First send identity with SNDMORE flag
	if _, err := socket.SendBytes(identity, zmq.SNDMORE); err != nil {
		return fmt.Errorf("failed to send identity: %w", err)
	}

	// Send data (final frame)
	if _, err := socket.SendBytes(data, 0); err != nil {
		return fmt.Errorf("failed to send data: %w", err)
	}

	return nil
}

// cleanupSession removes a session by identity
func (m *ZMQSessionManager) cleanupSession(identity []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	identityKey := string(identity)
	session, exists := m.identityToSession[identityKey]
	if !exists {
		return
	}

	// Remove session
	delete(m.sessions, session.ID)
	delete(m.identityToSession, identityKey)

	// Decrement session counts based on filter mask
	m.decrementSessionCounts(session.filterMask)

	log.Info("Session removed due to disconnect",
		"sessionID", session.ID,
		"activeSessions", m.sessionCounts[0].Load())
}

// shouldProcessEvent checks if event passes the filter
func (s *Session) shouldProcessEvent(event SBEEvent) bool {
	// Extract event type from SBE event
	var eventType uint8
	switch e := event.(type) {
	case *ethereum_tracing.TxStartEvent:
		eventType = uint8(e.EventType)
	case *ethereum_tracing.TxEndEvent:
		eventType = uint8(e.EventType)
	case *ethereum_tracing.CallEvent:
		eventType = uint8(e.EventType)
	case *ethereum_tracing.BlockStartEvent:
		eventType = uint8(e.EventType)
	case *ethereum_tracing.BlockEndEvent:
		eventType = uint8(e.EventType)
	case *ethereum_tracing.LogEvent:
		eventType = uint8(e.EventType)
		// LogEvent has additional filtering (eth_getLogs style)
		if !s.shouldProcessLogEvent(e) {
			return false
		}
	default:
		return false // Unknown event type
	}

	// Check event type filter
	if s.filterMask == 0 {
		return true // No filter, accept all
	}

	return (s.filterMask & (1 << eventType)) != 0
}

// shouldProcessLogEvent checks if LogEvent passes address and topic filters (eth_getLogs style)
func (s *Session) shouldProcessLogEvent(log *ethereum_tracing.LogEvent) bool {
	// Check address filter
	if len(s.logAddressFilter) > 0 {
		if !s.logAddressFilter[log.Address] {
			return false
		}
	}

	// Check topics filter (eth_getLogs style with wildcard support)
	if len(s.logTopicsFilter) > 0 {
		for i, topicFilter := range s.logTopicsFilter {
			// Topic filter position exceeds log topics count
			if i >= int(log.TopicsCount) {
				return false
			}

			// Empty map at this position = match any topic (wildcard)
			if len(topicFilter) == 0 {
				continue
			}

			// Check if log topic is in the filter map
			if !topicFilter[log.Topics[i].Topic] {
				return false
			}
		}
	}

	return true
}

// IsEmpty returns true if there are no active sessions
// This is a fast atomic check that can be called from hooks
func (m *ZMQSessionManager) IsEmpty() bool {
	return m.sessionCounts[0].Load() == 0
}

// SessionCount returns the current number of active sessions
func (m *ZMQSessionManager) SessionCount() int64 {
	return m.sessionCounts[0].Load()
}

// HasSessionsForEventType returns true if any session subscribes to the given event type
// This is a fast atomic check for hooks to skip event creation when no subscribers
func (m *ZMQSessionManager) HasSessionsForEventType(eventType ethereum_tracing.EventTypeEnum) bool {
	if eventType < ethereum_tracing.EventTypeEnum(m.eventTypeStart) || eventType >= ethereum_tracing.EventTypeEnum(m.eventTypeEnd) {
		return false
	}
	return m.sessionCounts[eventType].Load() > 0
}

// Stop stops the session manager
func (m *ZMQSessionManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return
	}

	m.closed = true

	// Signal shutdown first
	close(m.shutdownChan)

	// Wait for goroutines to exit gracefully (they will detect shutdown via timeout)
	// This avoids assertion failure in ZeroMQ when closing while RecvBytes is blocked
	m.wg.Wait()

	// Now safe to close all sockets (no goroutines using them)
	for _, socket := range m.routerSockets {
		socket.SetLinger(0)
		socket.Close()
	}

	log.Info("ZMQ tracer session manager stopped")
}

// encodeErrorResponse encodes an error response
func (m *ZMQSessionManager) encodeErrorResponse(errMsg string) []byte {
	resp := &ethereum_tracing.SessionCreateResponse{
		MessageType: ethereum_tracing.MessageType.SessionCreateResponse,
		Status:      ethereum_tracing.Status.Error,
		ErrorMsg:    []uint8(errMsg),
	}

	encoded, err := m.encoder.Encode(resp)
	if err != nil {
		// Fallback: encode minimal error response
		return []byte{2, 0, 0, 0, 0, 1} // SessionCreateResponse, 0, 0, 0, Error
	}
	return encoded
}

// encodeSuccessResponse encodes a success response
func (m *ZMQSessionManager) encodeSuccessResponse(sessionID uint32, endpoint string) []byte {
	resp := &ethereum_tracing.SessionCreateResponse{
		MessageType: ethereum_tracing.MessageType.SessionCreateResponse,
		SessionID:   sessionID,
		Status:      ethereum_tracing.Status.Success,
		PubEndpoint: []uint8(endpoint),
	}

	encoded, err := m.encoder.Encode(resp)
	if err != nil {
		// This should never happen, but fallback
		log.Error("Failed to encode success response", "error", err)
		return m.encodeErrorResponse("encoding error: " + err.Error())
	}
	return encoded
}
