// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"reflect"
)

type MessageTypeEnum uint8
type MessageTypeValues struct {
	SessionCreateRequest  MessageTypeEnum
	SessionCreateResponse MessageTypeEnum
	NullValue             MessageTypeEnum
}

var MessageType = MessageTypeValues{1, 2, 255}

func (m MessageTypeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteUint8(_w, uint8(m)); err != nil {
		return err
	}
	return nil
}

func (m *MessageTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadUint8(_r, (*uint8)(m)); err != nil {
		return err
	}
	return nil
}

func (m MessageTypeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(MessageType)
	for idx := 0; idx < value.NumField(); idx++ {
		if m == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on MessageType, unknown enumeration value %d", m)
}

func (*MessageTypeEnum) EncodedLength() int64 {
	return 1
}

func (*MessageTypeEnum) SessionCreateRequestSinceVersion() uint16 {
	return 0
}

func (m *MessageTypeEnum) SessionCreateRequestInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.SessionCreateRequestSinceVersion()
}

func (*MessageTypeEnum) SessionCreateRequestDeprecated() uint16 {
	return 0
}

func (*MessageTypeEnum) SessionCreateResponseSinceVersion() uint16 {
	return 0
}

func (m *MessageTypeEnum) SessionCreateResponseInActingVersion(actingVersion uint16) bool {
	return actingVersion >= m.SessionCreateResponseSinceVersion()
}

func (*MessageTypeEnum) SessionCreateResponseDeprecated() uint16 {
	return 0
}
