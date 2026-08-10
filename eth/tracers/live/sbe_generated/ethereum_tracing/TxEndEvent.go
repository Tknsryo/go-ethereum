// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TxEndEvent struct {
	EventType   EventTypeEnum
	Timestamp   uint64
	TxHash      [32]uint8
	Status      uint8
	GasUsed     uint64
	BlockNumber uint64
	BlockHash   [32]uint8
}

func (t *TxEndEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := t.RangeCheck(t.SbeSchemaVersion(), t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := t.EventType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, t.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.TxHash[:]); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, t.Status); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, t.GasUsed); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, t.BlockNumber); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.BlockHash[:]); err != nil {
		return err
	}
	return nil
}

func (t *TxEndEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if t.EventTypeInActingVersion(actingVersion) {
		if err := t.EventType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !t.TimestampInActingVersion(actingVersion) {
		t.Timestamp = t.TimestampNullValue()
	} else {
		if err := _m.ReadUint64(_r, &t.Timestamp); err != nil {
			return err
		}
	}
	if !t.TxHashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			t.TxHash[idx] = t.TxHashNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, t.TxHash[:]); err != nil {
			return err
		}
	}
	if !t.StatusInActingVersion(actingVersion) {
		t.Status = t.StatusNullValue()
	} else {
		if err := _m.ReadUint8(_r, &t.Status); err != nil {
			return err
		}
	}
	if !t.GasUsedInActingVersion(actingVersion) {
		t.GasUsed = t.GasUsedNullValue()
	} else {
		if err := _m.ReadUint64(_r, &t.GasUsed); err != nil {
			return err
		}
	}
	if !t.BlockNumberInActingVersion(actingVersion) {
		t.BlockNumber = t.BlockNumberNullValue()
	} else {
		if err := _m.ReadUint64(_r, &t.BlockNumber); err != nil {
			return err
		}
	}
	if !t.BlockHashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			t.BlockHash[idx] = t.BlockHashNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, t.BlockHash[:]); err != nil {
			return err
		}
	}
	if actingVersion > t.SbeSchemaVersion() && blockLength > t.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-t.SbeBlockLength()))
	}
	if doRangeCheck {
		if err := t.RangeCheck(actingVersion, t.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (t *TxEndEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if err := t.EventType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if t.TimestampInActingVersion(actingVersion) {
		if t.Timestamp < t.TimestampMinValue() || t.Timestamp > t.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on t.Timestamp (%v < %v > %v)", t.TimestampMinValue(), t.Timestamp, t.TimestampMaxValue())
		}
	}
	if t.TxHashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if t.TxHash[idx] < t.TxHashMinValue() || t.TxHash[idx] > t.TxHashMaxValue() {
				return fmt.Errorf("Range check failed on t.TxHash[%d] (%v < %v > %v)", idx, t.TxHashMinValue(), t.TxHash[idx], t.TxHashMaxValue())
			}
		}
	}
	if t.StatusInActingVersion(actingVersion) {
		if t.Status < t.StatusMinValue() || t.Status > t.StatusMaxValue() {
			return fmt.Errorf("Range check failed on t.Status (%v < %v > %v)", t.StatusMinValue(), t.Status, t.StatusMaxValue())
		}
	}
	if t.GasUsedInActingVersion(actingVersion) {
		if t.GasUsed < t.GasUsedMinValue() || t.GasUsed > t.GasUsedMaxValue() {
			return fmt.Errorf("Range check failed on t.GasUsed (%v < %v > %v)", t.GasUsedMinValue(), t.GasUsed, t.GasUsedMaxValue())
		}
	}
	if t.BlockNumberInActingVersion(actingVersion) {
		if t.BlockNumber < t.BlockNumberMinValue() || t.BlockNumber > t.BlockNumberMaxValue() {
			return fmt.Errorf("Range check failed on t.BlockNumber (%v < %v > %v)", t.BlockNumberMinValue(), t.BlockNumber, t.BlockNumberMaxValue())
		}
	}
	if t.BlockHashInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if t.BlockHash[idx] < t.BlockHashMinValue() || t.BlockHash[idx] > t.BlockHashMaxValue() {
				return fmt.Errorf("Range check failed on t.BlockHash[%d] (%v < %v > %v)", idx, t.BlockHashMinValue(), t.BlockHash[idx], t.BlockHashMaxValue())
			}
		}
	}
	return nil
}

func TxEndEventInit(t *TxEndEvent) {
	return
}

func (*TxEndEvent) SbeBlockLength() (blockLength uint16) {
	return 90
}

func (*TxEndEvent) SbeTemplateId() (templateId uint16) {
	return 11
}

func (*TxEndEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*TxEndEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*TxEndEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TxEndEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*TxEndEvent) EventTypeId() uint16 {
	return 1
}

func (*TxEndEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.EventTypeSinceVersion()
}

func (*TxEndEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) EventTypeMetaAttribute(meta int) string {
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

func (*TxEndEvent) TimestampId() uint16 {
	return 2
}

func (*TxEndEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.TimestampSinceVersion()
}

func (*TxEndEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) TimestampMetaAttribute(meta int) string {
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

func (*TxEndEvent) TimestampMinValue() uint64 {
	return 0
}

func (*TxEndEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TxEndEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*TxEndEvent) TxHashId() uint16 {
	return 4
}

func (*TxEndEvent) TxHashSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) TxHashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.TxHashSinceVersion()
}

func (*TxEndEvent) TxHashDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) TxHashMetaAttribute(meta int) string {
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

func (*TxEndEvent) TxHashMinValue() uint8 {
	return 0
}

func (*TxEndEvent) TxHashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxEndEvent) TxHashNullValue() uint8 {
	return math.MaxUint8
}

func (*TxEndEvent) StatusId() uint16 {
	return 5
}

func (*TxEndEvent) StatusSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) StatusInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.StatusSinceVersion()
}

func (*TxEndEvent) StatusDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) StatusMetaAttribute(meta int) string {
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

func (*TxEndEvent) StatusMinValue() uint8 {
	return 0
}

func (*TxEndEvent) StatusMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxEndEvent) StatusNullValue() uint8 {
	return math.MaxUint8
}

func (*TxEndEvent) GasUsedId() uint16 {
	return 6
}

func (*TxEndEvent) GasUsedSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) GasUsedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.GasUsedSinceVersion()
}

func (*TxEndEvent) GasUsedDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) GasUsedMetaAttribute(meta int) string {
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

func (*TxEndEvent) GasUsedMinValue() uint64 {
	return 0
}

func (*TxEndEvent) GasUsedMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TxEndEvent) GasUsedNullValue() uint64 {
	return math.MaxUint64
}

func (*TxEndEvent) BlockNumberId() uint16 {
	return 7
}

func (*TxEndEvent) BlockNumberSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) BlockNumberInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.BlockNumberSinceVersion()
}

func (*TxEndEvent) BlockNumberDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) BlockNumberMetaAttribute(meta int) string {
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

func (*TxEndEvent) BlockNumberMinValue() uint64 {
	return 0
}

func (*TxEndEvent) BlockNumberMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TxEndEvent) BlockNumberNullValue() uint64 {
	return math.MaxUint64
}

func (*TxEndEvent) BlockHashId() uint16 {
	return 8
}

func (*TxEndEvent) BlockHashSinceVersion() uint16 {
	return 0
}

func (t *TxEndEvent) BlockHashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.BlockHashSinceVersion()
}

func (*TxEndEvent) BlockHashDeprecated() uint16 {
	return 0
}

func (*TxEndEvent) BlockHashMetaAttribute(meta int) string {
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

func (*TxEndEvent) BlockHashMinValue() uint8 {
	return 0
}

func (*TxEndEvent) BlockHashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxEndEvent) BlockHashNullValue() uint8 {
	return math.MaxUint8
}
