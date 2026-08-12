// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type BlockEndEvent struct {
	EventType        EventTypeEnum
	Timestamp        uint64
	Number           uint64
	InsertDurationNs uint64
	ErrorMsg         []uint8
}

func (b *BlockEndEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := b.RangeCheck(b.SbeSchemaVersion(), b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := b.EventType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.Number); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.InsertDurationNs); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(b.ErrorMsg))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, b.ErrorMsg); err != nil {
		return err
	}
	return nil
}

func (b *BlockEndEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if b.EventTypeInActingVersion(actingVersion) {
		if err := b.EventType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !b.TimestampInActingVersion(actingVersion) {
		b.Timestamp = b.TimestampNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.Timestamp); err != nil {
			return err
		}
	}
	if !b.NumberInActingVersion(actingVersion) {
		b.Number = b.NumberNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.Number); err != nil {
			return err
		}
	}
	if !b.InsertDurationNsInActingVersion(actingVersion) {
		b.InsertDurationNs = b.InsertDurationNsNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.InsertDurationNs); err != nil {
			return err
		}
	}
	if actingVersion > b.SbeSchemaVersion() && blockLength > b.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}

	if b.ErrorMsgInActingVersion(actingVersion) {
		var ErrorMsgLength uint32
		if err := _m.ReadUint32(_r, &ErrorMsgLength); err != nil {
			return err
		}
		if cap(b.ErrorMsg) < int(ErrorMsgLength) {
			b.ErrorMsg = make([]uint8, ErrorMsgLength)
		}
		b.ErrorMsg = b.ErrorMsg[:ErrorMsgLength]
		if err := _m.ReadBytes(_r, b.ErrorMsg); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := b.RangeCheck(actingVersion, b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (b *BlockEndEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if err := b.EventType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if b.TimestampInActingVersion(actingVersion) {
		if b.Timestamp < b.TimestampMinValue() || b.Timestamp > b.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on b.Timestamp (%v < %v > %v)", b.TimestampMinValue(), b.Timestamp, b.TimestampMaxValue())
		}
	}
	if b.NumberInActingVersion(actingVersion) {
		if b.Number < b.NumberMinValue() || b.Number > b.NumberMaxValue() {
			return fmt.Errorf("Range check failed on b.Number (%v < %v > %v)", b.NumberMinValue(), b.Number, b.NumberMaxValue())
		}
	}
	if b.InsertDurationNsInActingVersion(actingVersion) {
		if b.InsertDurationNs < b.InsertDurationNsMinValue() || b.InsertDurationNs > b.InsertDurationNsMaxValue() {
			return fmt.Errorf("Range check failed on b.InsertDurationNs (%v < %v > %v)", b.InsertDurationNsMinValue(), b.InsertDurationNs, b.InsertDurationNsMaxValue())
		}
	}
	return nil
}

func BlockEndEventInit(b *BlockEndEvent) {
	return
}

func (*BlockEndEvent) SbeBlockLength() (blockLength uint16) {
	return 25
}

func (*BlockEndEvent) SbeTemplateId() (templateId uint16) {
	return 14
}

func (*BlockEndEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*BlockEndEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*BlockEndEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*BlockEndEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*BlockEndEvent) EventTypeId() uint16 {
	return 1
}

func (*BlockEndEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (b *BlockEndEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.EventTypeSinceVersion()
}

func (*BlockEndEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*BlockEndEvent) EventTypeMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*BlockEndEvent) TimestampId() uint16 {
	return 2
}

func (*BlockEndEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (b *BlockEndEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TimestampSinceVersion()
}

func (*BlockEndEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*BlockEndEvent) TimestampMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*BlockEndEvent) TimestampMinValue() uint64 {
	return 0
}

func (*BlockEndEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockEndEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockEndEvent) NumberId() uint16 {
	return 3
}

func (*BlockEndEvent) NumberSinceVersion() uint16 {
	return 0
}

func (b *BlockEndEvent) NumberInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.NumberSinceVersion()
}

func (*BlockEndEvent) NumberDeprecated() uint16 {
	return 0
}

func (*BlockEndEvent) NumberMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*BlockEndEvent) NumberMinValue() uint64 {
	return 0
}

func (*BlockEndEvent) NumberMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockEndEvent) NumberNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockEndEvent) InsertDurationNsId() uint16 {
	return 4
}

func (*BlockEndEvent) InsertDurationNsSinceVersion() uint16 {
	return 0
}

func (b *BlockEndEvent) InsertDurationNsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.InsertDurationNsSinceVersion()
}

func (*BlockEndEvent) InsertDurationNsDeprecated() uint16 {
	return 0
}

func (*BlockEndEvent) InsertDurationNsMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*BlockEndEvent) InsertDurationNsMinValue() uint64 {
	return 0
}

func (*BlockEndEvent) InsertDurationNsMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockEndEvent) InsertDurationNsNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockEndEvent) ErrorMsgMetaAttribute(meta int) string {
	switch meta {
	case 1:
		return ""
	case 2:
		return ""
	case 3:
		return ""
	case 4:
		return "required"
	}
	return ""
}

func (*BlockEndEvent) ErrorMsgSinceVersion() uint16 {
	return 0
}

func (b *BlockEndEvent) ErrorMsgInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ErrorMsgSinceVersion()
}

func (*BlockEndEvent) ErrorMsgDeprecated() uint16 {
	return 0
}

func (BlockEndEvent) ErrorMsgCharacterEncoding() string {
	return "null"
}

func (BlockEndEvent) ErrorMsgHeaderLength() uint64 {
	return 4
}
