// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type BlockStartEvent struct {
	EventType      EventTypeEnum
	Timestamp      uint64
	Number         uint64
	Hash           [32]uint8
	ParentHash     [32]uint8
	BlockTimestamp uint64
	TxCount        uint16
	GasUsed        uint64
	GasLimit       uint64
	BaseFeePerGas  uint64
}

func (b *BlockStartEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteUint16(_w, b.TxCount); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.GasUsed); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.GasLimit); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, b.BaseFeePerGas); err != nil {
		return err
	}
	return nil
}

func (b *BlockStartEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !b.TxCountInActingVersion(actingVersion) {
		b.TxCount = b.TxCountNullValue()
	} else {
		if err := _m.ReadUint16(_r, &b.TxCount); err != nil {
			return err
		}
	}
	if !b.GasUsedInActingVersion(actingVersion) {
		b.GasUsed = b.GasUsedNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.GasUsed); err != nil {
			return err
		}
	}
	if !b.GasLimitInActingVersion(actingVersion) {
		b.GasLimit = b.GasLimitNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.GasLimit); err != nil {
			return err
		}
	}
	if !b.BaseFeePerGasInActingVersion(actingVersion) {
		b.BaseFeePerGas = b.BaseFeePerGasNullValue()
	} else {
		if err := _m.ReadUint64(_r, &b.BaseFeePerGas); err != nil {
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

func (b *BlockStartEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if b.TxCountInActingVersion(actingVersion) {
		if b.TxCount < b.TxCountMinValue() || b.TxCount > b.TxCountMaxValue() {
			return fmt.Errorf("Range check failed on b.TxCount (%v < %v > %v)", b.TxCountMinValue(), b.TxCount, b.TxCountMaxValue())
		}
	}
	if b.GasUsedInActingVersion(actingVersion) {
		if b.GasUsed < b.GasUsedMinValue() || b.GasUsed > b.GasUsedMaxValue() {
			return fmt.Errorf("Range check failed on b.GasUsed (%v < %v > %v)", b.GasUsedMinValue(), b.GasUsed, b.GasUsedMaxValue())
		}
	}
	if b.GasLimitInActingVersion(actingVersion) {
		if b.GasLimit < b.GasLimitMinValue() || b.GasLimit > b.GasLimitMaxValue() {
			return fmt.Errorf("Range check failed on b.GasLimit (%v < %v > %v)", b.GasLimitMinValue(), b.GasLimit, b.GasLimitMaxValue())
		}
	}
	if b.BaseFeePerGasInActingVersion(actingVersion) {
		if b.BaseFeePerGas < b.BaseFeePerGasMinValue() || b.BaseFeePerGas > b.BaseFeePerGasMaxValue() {
			return fmt.Errorf("Range check failed on b.BaseFeePerGas (%v < %v > %v)", b.BaseFeePerGasMinValue(), b.BaseFeePerGas, b.BaseFeePerGasMaxValue())
		}
	}
	return nil
}

func BlockStartEventInit(b *BlockStartEvent) {
	return
}

func (*BlockStartEvent) SbeBlockLength() (blockLength uint16) {
	return 115
}

func (*BlockStartEvent) SbeTemplateId() (templateId uint16) {
	return 13
}

func (*BlockStartEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*BlockStartEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*BlockStartEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*BlockStartEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*BlockStartEvent) EventTypeId() uint16 {
	return 1
}

func (*BlockStartEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.EventTypeSinceVersion()
}

func (*BlockStartEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) EventTypeMetaAttribute(meta int) string {
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

func (*BlockStartEvent) TimestampId() uint16 {
	return 2
}

func (*BlockStartEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TimestampSinceVersion()
}

func (*BlockStartEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) TimestampMetaAttribute(meta int) string {
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

func (*BlockStartEvent) TimestampMinValue() uint64 {
	return 0
}

func (*BlockStartEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockStartEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockStartEvent) NumberId() uint16 {
	return 3
}

func (*BlockStartEvent) NumberSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) NumberInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.NumberSinceVersion()
}

func (*BlockStartEvent) NumberDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) NumberMetaAttribute(meta int) string {
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

func (*BlockStartEvent) NumberMinValue() uint64 {
	return 0
}

func (*BlockStartEvent) NumberMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockStartEvent) NumberNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockStartEvent) HashId() uint16 {
	return 4
}

func (*BlockStartEvent) HashSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) HashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.HashSinceVersion()
}

func (*BlockStartEvent) HashDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) HashMetaAttribute(meta int) string {
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

func (*BlockStartEvent) HashMinValue() uint8 {
	return 0
}

func (*BlockStartEvent) HashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*BlockStartEvent) HashNullValue() uint8 {
	return math.MaxUint8
}

func (*BlockStartEvent) ParentHashId() uint16 {
	return 5
}

func (*BlockStartEvent) ParentHashSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) ParentHashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.ParentHashSinceVersion()
}

func (*BlockStartEvent) ParentHashDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) ParentHashMetaAttribute(meta int) string {
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

func (*BlockStartEvent) ParentHashMinValue() uint8 {
	return 0
}

func (*BlockStartEvent) ParentHashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*BlockStartEvent) ParentHashNullValue() uint8 {
	return math.MaxUint8
}

func (*BlockStartEvent) BlockTimestampId() uint16 {
	return 6
}

func (*BlockStartEvent) BlockTimestampSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) BlockTimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.BlockTimestampSinceVersion()
}

func (*BlockStartEvent) BlockTimestampDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) BlockTimestampMetaAttribute(meta int) string {
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

func (*BlockStartEvent) BlockTimestampMinValue() uint64 {
	return 0
}

func (*BlockStartEvent) BlockTimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockStartEvent) BlockTimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockStartEvent) TxCountId() uint16 {
	return 7
}

func (*BlockStartEvent) TxCountSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) TxCountInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.TxCountSinceVersion()
}

func (*BlockStartEvent) TxCountDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) TxCountMetaAttribute(meta int) string {
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

func (*BlockStartEvent) TxCountMinValue() uint16 {
	return 0
}

func (*BlockStartEvent) TxCountMaxValue() uint16 {
	return math.MaxUint16 - 1
}

func (*BlockStartEvent) TxCountNullValue() uint16 {
	return math.MaxUint16
}

func (*BlockStartEvent) GasUsedId() uint16 {
	return 8
}

func (*BlockStartEvent) GasUsedSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) GasUsedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.GasUsedSinceVersion()
}

func (*BlockStartEvent) GasUsedDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) GasUsedMetaAttribute(meta int) string {
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

func (*BlockStartEvent) GasUsedMinValue() uint64 {
	return 0
}

func (*BlockStartEvent) GasUsedMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockStartEvent) GasUsedNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockStartEvent) GasLimitId() uint16 {
	return 9
}

func (*BlockStartEvent) GasLimitSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) GasLimitInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.GasLimitSinceVersion()
}

func (*BlockStartEvent) GasLimitDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) GasLimitMetaAttribute(meta int) string {
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

func (*BlockStartEvent) GasLimitMinValue() uint64 {
	return 0
}

func (*BlockStartEvent) GasLimitMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockStartEvent) GasLimitNullValue() uint64 {
	return math.MaxUint64
}

func (*BlockStartEvent) BaseFeePerGasId() uint16 {
	return 10
}

func (*BlockStartEvent) BaseFeePerGasSinceVersion() uint16 {
	return 0
}

func (b *BlockStartEvent) BaseFeePerGasInActingVersion(actingVersion uint16) bool {
	return actingVersion >= b.BaseFeePerGasSinceVersion()
}

func (*BlockStartEvent) BaseFeePerGasDeprecated() uint16 {
	return 0
}

func (*BlockStartEvent) BaseFeePerGasMetaAttribute(meta int) string {
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

func (*BlockStartEvent) BaseFeePerGasMinValue() uint64 {
	return 0
}

func (*BlockStartEvent) BaseFeePerGasMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*BlockStartEvent) BaseFeePerGasNullValue() uint64 {
	return math.MaxUint64
}
