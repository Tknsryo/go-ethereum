// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type TxStartEvent struct {
	EventType EventTypeEnum
	Timestamp uint64
	TxHash    [32]uint8
	From      [20]uint8
	To        [20]uint8
	Value     [32]uint8
	GasLimit  uint64
	GasPrice  [32]uint8
	Nonce     uint64
}

func (t *TxStartEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
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
	if err := _m.WriteBytes(_w, t.From[:]); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.To[:]); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.Value[:]); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, t.GasLimit); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, t.GasPrice[:]); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, t.Nonce); err != nil {
		return err
	}
	return nil
}

func (t *TxStartEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
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
	if !t.FromInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			t.From[idx] = t.FromNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, t.From[:]); err != nil {
			return err
		}
	}
	if !t.ToInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			t.To[idx] = t.ToNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, t.To[:]); err != nil {
			return err
		}
	}
	if !t.ValueInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			t.Value[idx] = t.ValueNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, t.Value[:]); err != nil {
			return err
		}
	}
	if !t.GasLimitInActingVersion(actingVersion) {
		t.GasLimit = t.GasLimitNullValue()
	} else {
		if err := _m.ReadUint64(_r, &t.GasLimit); err != nil {
			return err
		}
	}
	if !t.GasPriceInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			t.GasPrice[idx] = t.GasPriceNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, t.GasPrice[:]); err != nil {
			return err
		}
	}
	if !t.NonceInActingVersion(actingVersion) {
		t.Nonce = t.NonceNullValue()
	} else {
		if err := _m.ReadUint64(_r, &t.Nonce); err != nil {
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

func (t *TxStartEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
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
	if t.FromInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			if t.From[idx] < t.FromMinValue() || t.From[idx] > t.FromMaxValue() {
				return fmt.Errorf("Range check failed on t.From[%d] (%v < %v > %v)", idx, t.FromMinValue(), t.From[idx], t.FromMaxValue())
			}
		}
	}
	if t.ToInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			if t.To[idx] < t.ToMinValue() || t.To[idx] > t.ToMaxValue() {
				return fmt.Errorf("Range check failed on t.To[%d] (%v < %v > %v)", idx, t.ToMinValue(), t.To[idx], t.ToMaxValue())
			}
		}
	}
	if t.ValueInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if t.Value[idx] < t.ValueMinValue() || t.Value[idx] > t.ValueMaxValue() {
				return fmt.Errorf("Range check failed on t.Value[%d] (%v < %v > %v)", idx, t.ValueMinValue(), t.Value[idx], t.ValueMaxValue())
			}
		}
	}
	if t.GasLimitInActingVersion(actingVersion) {
		if t.GasLimit < t.GasLimitMinValue() || t.GasLimit > t.GasLimitMaxValue() {
			return fmt.Errorf("Range check failed on t.GasLimit (%v < %v > %v)", t.GasLimitMinValue(), t.GasLimit, t.GasLimitMaxValue())
		}
	}
	if t.GasPriceInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if t.GasPrice[idx] < t.GasPriceMinValue() || t.GasPrice[idx] > t.GasPriceMaxValue() {
				return fmt.Errorf("Range check failed on t.GasPrice[%d] (%v < %v > %v)", idx, t.GasPriceMinValue(), t.GasPrice[idx], t.GasPriceMaxValue())
			}
		}
	}
	if t.NonceInActingVersion(actingVersion) {
		if t.Nonce < t.NonceMinValue() || t.Nonce > t.NonceMaxValue() {
			return fmt.Errorf("Range check failed on t.Nonce (%v < %v > %v)", t.NonceMinValue(), t.Nonce, t.NonceMaxValue())
		}
	}
	return nil
}

func TxStartEventInit(t *TxStartEvent) {
	return
}

func (*TxStartEvent) SbeBlockLength() (blockLength uint16) {
	return 161
}

func (*TxStartEvent) SbeTemplateId() (templateId uint16) {
	return 10
}

func (*TxStartEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*TxStartEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*TxStartEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*TxStartEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*TxStartEvent) EventTypeId() uint16 {
	return 1
}

func (*TxStartEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.EventTypeSinceVersion()
}

func (*TxStartEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) EventTypeMetaAttribute(meta int) string {
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

func (*TxStartEvent) TimestampId() uint16 {
	return 2
}

func (*TxStartEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.TimestampSinceVersion()
}

func (*TxStartEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) TimestampMetaAttribute(meta int) string {
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

func (*TxStartEvent) TimestampMinValue() uint64 {
	return 0
}

func (*TxStartEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TxStartEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*TxStartEvent) TxHashId() uint16 {
	return 4
}

func (*TxStartEvent) TxHashSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) TxHashInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.TxHashSinceVersion()
}

func (*TxStartEvent) TxHashDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) TxHashMetaAttribute(meta int) string {
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

func (*TxStartEvent) TxHashMinValue() uint8 {
	return 0
}

func (*TxStartEvent) TxHashMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxStartEvent) TxHashNullValue() uint8 {
	return math.MaxUint8
}

func (*TxStartEvent) FromId() uint16 {
	return 5
}

func (*TxStartEvent) FromSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) FromInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.FromSinceVersion()
}

func (*TxStartEvent) FromDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) FromMetaAttribute(meta int) string {
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

func (*TxStartEvent) FromMinValue() uint8 {
	return 0
}

func (*TxStartEvent) FromMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxStartEvent) FromNullValue() uint8 {
	return math.MaxUint8
}

func (*TxStartEvent) ToId() uint16 {
	return 6
}

func (*TxStartEvent) ToSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) ToInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.ToSinceVersion()
}

func (*TxStartEvent) ToDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) ToMetaAttribute(meta int) string {
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

func (*TxStartEvent) ToMinValue() uint8 {
	return 0
}

func (*TxStartEvent) ToMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxStartEvent) ToNullValue() uint8 {
	return math.MaxUint8
}

func (*TxStartEvent) ValueId() uint16 {
	return 7
}

func (*TxStartEvent) ValueSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) ValueInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.ValueSinceVersion()
}

func (*TxStartEvent) ValueDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) ValueMetaAttribute(meta int) string {
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

func (*TxStartEvent) ValueMinValue() uint8 {
	return 0
}

func (*TxStartEvent) ValueMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxStartEvent) ValueNullValue() uint8 {
	return math.MaxUint8
}

func (*TxStartEvent) GasLimitId() uint16 {
	return 8
}

func (*TxStartEvent) GasLimitSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) GasLimitInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.GasLimitSinceVersion()
}

func (*TxStartEvent) GasLimitDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) GasLimitMetaAttribute(meta int) string {
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

func (*TxStartEvent) GasLimitMinValue() uint64 {
	return 0
}

func (*TxStartEvent) GasLimitMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TxStartEvent) GasLimitNullValue() uint64 {
	return math.MaxUint64
}

func (*TxStartEvent) GasPriceId() uint16 {
	return 9
}

func (*TxStartEvent) GasPriceSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) GasPriceInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.GasPriceSinceVersion()
}

func (*TxStartEvent) GasPriceDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) GasPriceMetaAttribute(meta int) string {
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

func (*TxStartEvent) GasPriceMinValue() uint8 {
	return 0
}

func (*TxStartEvent) GasPriceMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*TxStartEvent) GasPriceNullValue() uint8 {
	return math.MaxUint8
}

func (*TxStartEvent) NonceId() uint16 {
	return 10
}

func (*TxStartEvent) NonceSinceVersion() uint16 {
	return 0
}

func (t *TxStartEvent) NonceInActingVersion(actingVersion uint16) bool {
	return actingVersion >= t.NonceSinceVersion()
}

func (*TxStartEvent) NonceDeprecated() uint16 {
	return 0
}

func (*TxStartEvent) NonceMetaAttribute(meta int) string {
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

func (*TxStartEvent) NonceMinValue() uint64 {
	return 0
}

func (*TxStartEvent) NonceMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*TxStartEvent) NonceNullValue() uint64 {
	return math.MaxUint64
}
