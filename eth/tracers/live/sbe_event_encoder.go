package live

import (
	"bytes"
	"io"

	"github.com/ethereum/go-ethereum/eth/tracers/live/sbe_generated/ethereum_tracing"
)

// SBEEvent is an interface for all SBE-generated event types
type SBEEvent interface {
	Encode(_m *ethereum_tracing.SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error
	SbeTemplateId() uint16
	SbeBlockLength() uint16
	SbeSchemaId() uint16
	SbeSchemaVersion() uint16
}

// SBEEventEncoder encodes SBE events directly using generated code
type SBEEventEncoder struct {
	marshaller *ethereum_tracing.SbeGoMarshaller
	buffer     *bytes.Buffer
}

// NewSBEEventEncoder creates a new encoder
func NewSBEEventEncoder() *SBEEventEncoder {
	return &SBEEventEncoder{
		marshaller: ethereum_tracing.NewSbeGoMarshaller(),
		buffer:     bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

// Encode encodes an SBE event with header
func (e *SBEEventEncoder) Encode(event SBEEvent) ([]byte, error) {
	e.buffer.Reset()

	// Encode message header
	header := ethereum_tracing.SbeGoMessageHeader{
		BlockLength: event.SbeBlockLength(),
		TemplateId:  event.SbeTemplateId(),
		SchemaId:    event.SbeSchemaId(),
		Version:     event.SbeSchemaVersion(),
	}

	if err := header.Encode(e.marshaller, e.buffer); err != nil {
		return nil, err
	}

	// Encode message body (skip range check for performance)
	if err := event.Encode(e.marshaller, e.buffer, false); err != nil {
		return nil, err
	}

	return e.buffer.Bytes(), nil
}