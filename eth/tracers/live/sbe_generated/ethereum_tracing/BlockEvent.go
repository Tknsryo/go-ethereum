// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type BlockEvent struct {
	EventType      EventTypeEnum
	Timestamp      uint64
	Number         uint64
	Hash           [32]uint8
	ParentHash     [32]uint8
	BlockTimestamp uint64
}

func (b *BlockEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteBytes(_w, b.Hash[:]); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, b.ParentHash[:]); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.BlockTimestamp); err != nil {
		return err
	}
	return nil
}

func (b *BlockEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !b.HashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			b.Hash[idx] = b.HashNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, b.Hash[:]); err != nil {
			return err
		}
	}
	if !b.ParentHashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			b.ParentHash[idx] = b.ParentHashNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, b.ParentHash[:]); err != nil {
			return err
		}
	}
	if !b.BlockTimestampInActingVersion(actingVersion) {
		b.BlockTimestamp = b.BlockTimestampNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.BlockTimestamp); err != nil {
			return err
		}
	}
	if actingVersion > b.SbeSchemaVersion() && blockLength > b.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-b.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := b.RangeCheck(actingVersion, b.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (b *BlockEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if b.HashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if b.Hash[idx] < b.HashMinValue() || b.Hash[idx] > b.HashMaxValue() {
				return fmt.Errorf("Range check failed on b.Hash[%d] (%v < %v > %v)", idx, b.HashMinValue(), b.Hash[idx], b.HashMaxValue())
			}
		}
	}
	if b.ParentHashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if b.ParentHash[idx] < b.ParentHashMinValue() || b.ParentHash[idx] > b.ParentHashMaxValue() {
				return fmt.Errorf("Range check failed on b.ParentHash[%d] (%v < %v > %v)", idx, b.ParentHashMinValue(), b.ParentHash[idx], b.ParentHashMaxValue())
			}
		}
	}
	if b.BlockTimestampInActingVersion(actingVersion) {
		if b.BlockTimestamp < b.BlockTimestampMinValue() || b.BlockTimestamp > b.BlockTimestampMaxValue() {
			return fmt.Errorf("Range check failed on b.BlockTimestamp (%v < %v > %v)", b.BlockTimestampMinValue(), b.BlockTimestamp, b.BlockTimestampMaxValue())
		}
	}
	return nil
}

func BlockEventInit(b *BlockEvent) {
	return
}

func (*BlockEvent) SbeBlockLength() (blockLength uint16) {
	return 89
}

func (*BlockEvent) SbeTemplateId() (templateId uint16) {
	return 13
}

func (*BlockEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*BlockEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*BlockEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*BlockEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*BlockEvent) EventTypeId() uint16 {
	return 1
}

func (*BlockEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (b *BlockEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.EventTypeSinceVersion()
}

func (*BlockEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*BlockEvent) EventTypeMetaAttribute(meta int) string {
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

func (*BlockEvent) TimestampId() uint16 {
	return 2
}

func (*BlockEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (b *BlockEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TimestampSinceVersion()
}

func (*BlockEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*BlockEvent) TimestampMetaAttribute(meta int) string {
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

func (*BlockEvent) TimestampMinValue() uint64 {
	return 0
}

func (*BlockEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockEvent) NumberId() uint16 {
	return 4
}

func (*BlockEvent) NumberSinceVersion() uint16 {
	return 0
}

func (b *BlockEvent) NumberInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.NumberSinceVersion()
}

func (*BlockEvent) NumberDeprecated() uint16 {
	return 0
}

func (*BlockEvent) NumberMetaAttribute(meta int) string {
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

func (*BlockEvent) NumberMinValue() uint64 {
	return 0
}

func (*BlockEvent) NumberMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockEvent) NumberNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockEvent) HashId() uint16 {
	return 5
}

func (*BlockEvent) HashSinceVersion() uint16 {
	return 0
}

func (b *BlockEvent) HashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.HashSinceVersion()
}

func (*BlockEvent) HashDeprecated() uint16 {
	return 0
}

func (*BlockEvent) HashMetaAttribute(meta int) string {
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

func (*BlockEvent) HashMinValue() uint8 {
	return 0
}

func (*BlockEvent) HashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*BlockEvent) HashNullValue() uint8 {
	return math.MaxUint8
}

func (*BlockEvent) ParentHashId() uint16 {
	return 6
}

func (*BlockEvent) ParentHashSinceVersion() uint16 {
	return 0
}

func (b *BlockEvent) ParentHashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ParentHashSinceVersion()
}

func (*BlockEvent) ParentHashDeprecated() uint16 {
	return 0
}

func (*BlockEvent) ParentHashMetaAttribute(meta int) string {
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

func (*BlockEvent) ParentHashMinValue() uint8 {
	return 0
}

func (*BlockEvent) ParentHashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*BlockEvent) ParentHashNullValue() uint8 {
	return math.MaxUint8
}

func (*BlockEvent) BlockTimestampId() uint16 {
	return 7
}

func (*BlockEvent) BlockTimestampSinceVersion() uint16 {
	return 0
}

func (b *BlockEvent) BlockTimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.BlockTimestampSinceVersion()
}

func (*BlockEvent) BlockTimestampDeprecated() uint16 {
	return 0
}

func (*BlockEvent) BlockTimestampMetaAttribute(meta int) string {
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

func (*BlockEvent) BlockTimestampMinValue() uint64 {
	return 0
}

func (*BlockEvent) BlockTimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockEvent) BlockTimestampNullValue() uint64 {
	return math.MaxUint64
}
