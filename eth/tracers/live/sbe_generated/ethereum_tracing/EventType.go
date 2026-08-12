// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"reflect"
)

type EventTypeEnum uint8
type EventTypeValues struct {
	TxStart    EventTypeEnum
	TxEnd      EventTypeEnum
	Enter      EventTypeEnum
	Exit       EventTypeEnum
	BlockStart EventTypeEnum
	BlockEnd   EventTypeEnum
	Log        EventTypeEnum
	NullValue  EventTypeEnum
}

var EventType = EventTypeValues{1, 2, 3, 4, 5, 6, 7, 255}

func (e EventTypeEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteUint8(_w, uint8(e)); err != nil {
		return err
	}
	return nil
}

func (e *EventTypeEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadUint8(_r, (*uint8)(e)); err != nil {
		return err
	}
	return nil
}

func (e EventTypeEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(EventType)
	for idx := 0; idx < value.NumField(); idx++ {
		if e == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on EventType, unknown enumeration value %d", e)
}

func (*EventTypeEnum) EncodedLength() int64 {
	return 1
}

func (*EventTypeEnum) TxStartSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) TxStartInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.TxStartSinceVersion()
}

func (*EventTypeEnum) TxStartDeprecated() uint16 {
	return 0
}

func (*EventTypeEnum) TxEndSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) TxEndInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.TxEndSinceVersion()
}

func (*EventTypeEnum) TxEndDeprecated() uint16 {
	return 0
}

func (*EventTypeEnum) EnterSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) EnterInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.EnterSinceVersion()
}

func (*EventTypeEnum) EnterDeprecated() uint16 {
	return 0
}

func (*EventTypeEnum) ExitSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) ExitInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.ExitSinceVersion()
}

func (*EventTypeEnum) ExitDeprecated() uint16 {
	return 0
}

func (*EventTypeEnum) BlockStartSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) BlockStartInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.BlockStartSinceVersion()
}

func (*EventTypeEnum) BlockStartDeprecated() uint16 {
	return 0
}

func (*EventTypeEnum) BlockEndSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) BlockEndInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.BlockEndSinceVersion()
}

func (*EventTypeEnum) BlockEndDeprecated() uint16 {
	return 0
}

func (*EventTypeEnum) LogSinceVersion() uint16 {
	return 0
}

func (e *EventTypeEnum) LogInActingVersion(actingVersion uint16) bool {
	return actingVersion >= e.LogSinceVersion()
}

func (*EventTypeEnum) LogDeprecated() uint16 {
	return 0
}
