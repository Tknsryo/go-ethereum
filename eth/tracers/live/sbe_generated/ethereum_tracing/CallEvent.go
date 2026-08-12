// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type CallEvent struct {
	EventType  EventTypeEnum
	Timestamp  uint64
	Depth      uint8
	CallType   uint8
	From       [20]uint8
	To         [20]uint8
	Gas        uint64
	Value      [32]uint8
	OutputSize uint32
	Reverted   uint8
	Input      []uint8
	Output     []uint8
}

func (c *CallEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := c.RangeCheck(c.SbeSchemaVersion(), c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := c.EventType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, c.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, c.Depth); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, c.CallType); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.From[:]); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.To[:]); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, c.Gas); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.Value[:]); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, c.OutputSize); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, c.Reverted); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.Input))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.Input); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(c.Output))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, c.Output); err != nil {
		return err
	}
	return nil
}

func (c *CallEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if c.EventTypeInActingVersion(actingVersion) {
		if err := c.EventType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !c.TimestampInActingVersion(actingVersion) {
		c.Timestamp = c.TimestampNullValue()
	} else {
		if err := _m.ReadUint64(_r, &c.Timestamp); err != nil {
			return err
		}
	}
	if !c.DepthInActingVersion(actingVersion) {
		c.Depth = c.DepthNullValue()
	} else {
		if err := _m.ReadUint8(_r, &c.Depth); err != nil {
			return err
		}
	}
	if !c.CallTypeInActingVersion(actingVersion) {
		c.CallType = c.CallTypeNullValue()
	} else {
		if err := _m.ReadUint8(_r, &c.CallType); err != nil {
			return err
		}
	}
	if !c.FromInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			c.From[idx] = c.FromNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, c.From[:]); err != nil {
			return err
		}
	}
	if !c.ToInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			c.To[idx] = c.ToNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, c.To[:]); err != nil {
			return err
		}
	}
	if !c.GasInActingVersion(actingVersion) {
		c.Gas = c.GasNullValue()
	} else {
		if err := _m.ReadUint64(_r, &c.Gas); err != nil {
			return err
		}
	}
	if !c.ValueInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			c.Value[idx] = c.ValueNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, c.Value[:]); err != nil {
			return err
		}
	}
	if !c.OutputSizeInActingVersion(actingVersion) {
		c.OutputSize = c.OutputSizeNullValue()
	} else {
		if err := _m.ReadUint32(_r, &c.OutputSize); err != nil {
			return err
		}
	}
	if !c.RevertedInActingVersion(actingVersion) {
		c.Reverted = c.RevertedNullValue()
	} else {
		if err := _m.ReadUint8(_r, &c.Reverted); err != nil {
			return err
		}
	}
	if actingVersion > c.SbeSchemaVersion() && blockLength > c.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-c.SbeBlockLength()))
	}

	if c.InputInActingVersion(actingVersion) {
		var InputLength uint32
		if err := _m.ReadUint32(_r, &InputLength); err != nil {
			return err
		}
		if cap(c.Input) < int(InputLength) {
			c.Input = make([]uint8, InputLength)
		}
		c.Input = c.Input[:InputLength]
		if err := _m.ReadBytes(_r, c.Input); err != nil {
			return err
		}
	}

	if c.OutputInActingVersion(actingVersion) {
		var OutputLength uint32
		if err := _m.ReadUint32(_r, &OutputLength); err != nil {
			return err
		}
		if cap(c.Output) < int(OutputLength) {
			c.Output = make([]uint8, OutputLength)
		}
		c.Output = c.Output[:OutputLength]
		if err := _m.ReadBytes(_r, c.Output); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := c.RangeCheck(actingVersion, c.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (c *CallEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if err := c.EventType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if c.TimestampInActingVersion(actingVersion) {
		if c.Timestamp < c.TimestampMinValue() || c.Timestamp > c.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on c.Timestamp (%v < %v > %v)", c.TimestampMinValue(), c.Timestamp, c.TimestampMaxValue())
		}
	}
	if c.DepthInActingVersion(actingVersion) {
		if c.Depth < c.DepthMinValue() || c.Depth > c.DepthMaxValue() {
			return fmt.Errorf("Range check failed on c.Depth (%v < %v > %v)", c.DepthMinValue(), c.Depth, c.DepthMaxValue())
		}
	}
	if c.CallTypeInActingVersion(actingVersion) {
		if c.CallType < c.CallTypeMinValue() || c.CallType > c.CallTypeMaxValue() {
			return fmt.Errorf("Range check failed on c.CallType (%v < %v > %v)", c.CallTypeMinValue(), c.CallType, c.CallTypeMaxValue())
		}
	}
	if c.FromInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			if c.From[idx] < c.FromMinValue() || c.From[idx] > c.FromMaxValue() {
				return fmt.Errorf("Range check failed on c.From[%d] (%v < %v > %v)", idx, c.FromMinValue(), c.From[idx], c.FromMaxValue())
			}
		}
	}
	if c.ToInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			if c.To[idx] < c.ToMinValue() || c.To[idx] > c.ToMaxValue() {
				return fmt.Errorf("Range check failed on c.To[%d] (%v < %v > %v)", idx, c.ToMinValue(), c.To[idx], c.ToMaxValue())
			}
		}
	}
	if c.GasInActingVersion(actingVersion) {
		if c.Gas < c.GasMinValue() || c.Gas > c.GasMaxValue() {
			return fmt.Errorf("Range check failed on c.Gas (%v < %v > %v)", c.GasMinValue(), c.Gas, c.GasMaxValue())
		}
	}
	if c.ValueInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if c.Value[idx] < c.ValueMinValue() || c.Value[idx] > c.ValueMaxValue() {
				return fmt.Errorf("Range check failed on c.Value[%d] (%v < %v > %v)", idx, c.ValueMinValue(), c.Value[idx], c.ValueMaxValue())
			}
		}
	}
	if c.OutputSizeInActingVersion(actingVersion) {
		if c.OutputSize < c.OutputSizeMinValue() || c.OutputSize > c.OutputSizeMaxValue() {
			return fmt.Errorf("Range check failed on c.OutputSize (%v < %v > %v)", c.OutputSizeMinValue(), c.OutputSize, c.OutputSizeMaxValue())
		}
	}
	if c.RevertedInActingVersion(actingVersion) {
		if c.Reverted < c.RevertedMinValue() || c.Reverted > c.RevertedMaxValue() {
			return fmt.Errorf("Range check failed on c.Reverted (%v < %v > %v)", c.RevertedMinValue(), c.Reverted, c.RevertedMaxValue())
		}
	}
	return nil
}

func CallEventInit(c *CallEvent) {
	return
}

func (*CallEvent) SbeBlockLength() (blockLength uint16) {
	return 96
}

func (*CallEvent) SbeTemplateId() (templateId uint16) {
	return 12
}

func (*CallEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*CallEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*CallEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*CallEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*CallEvent) EventTypeId() uint16 {
	return 1
}

func (*CallEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.EventTypeSinceVersion()
}

func (*CallEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*CallEvent) EventTypeMetaAttribute(meta int) string {
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

func (*CallEvent) TimestampId() uint16 {
	return 2
}

func (*CallEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.TimestampSinceVersion()
}

func (*CallEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*CallEvent) TimestampMetaAttribute(meta int) string {
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

func (*CallEvent) TimestampMinValue() uint64 {
	return 0
}

func (*CallEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*CallEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*CallEvent) DepthId() uint16 {
	return 4
}

func (*CallEvent) DepthSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) DepthInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.DepthSinceVersion()
}

func (*CallEvent) DepthDeprecated() uint16 {
	return 0
}

func (*CallEvent) DepthMetaAttribute(meta int) string {
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

func (*CallEvent) DepthMinValue() uint8 {
	return 0
}

func (*CallEvent) DepthMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*CallEvent) DepthNullValue() uint8 {
	return math.MaxUint8
}

func (*CallEvent) CallTypeId() uint16 {
	return 5
}

func (*CallEvent) CallTypeSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) CallTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.CallTypeSinceVersion()
}

func (*CallEvent) CallTypeDeprecated() uint16 {
	return 0
}

func (*CallEvent) CallTypeMetaAttribute(meta int) string {
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

func (*CallEvent) CallTypeMinValue() uint8 {
	return 0
}

func (*CallEvent) CallTypeMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*CallEvent) CallTypeNullValue() uint8 {
	return math.MaxUint8
}

func (*CallEvent) FromId() uint16 {
	return 6
}

func (*CallEvent) FromSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) FromInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.FromSinceVersion()
}

func (*CallEvent) FromDeprecated() uint16 {
	return 0
}

func (*CallEvent) FromMetaAttribute(meta int) string {
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

func (*CallEvent) FromMinValue() uint8 {
	return 0
}

func (*CallEvent) FromMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*CallEvent) FromNullValue() uint8 {
	return math.MaxUint8
}

func (*CallEvent) ToId() uint16 {
	return 7
}

func (*CallEvent) ToSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) ToInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ToSinceVersion()
}

func (*CallEvent) ToDeprecated() uint16 {
	return 0
}

func (*CallEvent) ToMetaAttribute(meta int) string {
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

func (*CallEvent) ToMinValue() uint8 {
	return 0
}

func (*CallEvent) ToMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*CallEvent) ToNullValue() uint8 {
	return math.MaxUint8
}

func (*CallEvent) GasId() uint16 {
	return 8
}

func (*CallEvent) GasSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) GasInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.GasSinceVersion()
}

func (*CallEvent) GasDeprecated() uint16 {
	return 0
}

func (*CallEvent) GasMetaAttribute(meta int) string {
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

func (*CallEvent) GasMinValue() uint64 {
	return 0
}

func (*CallEvent) GasMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*CallEvent) GasNullValue() uint64 {
	return math.MaxUint64
}

func (*CallEvent) ValueId() uint16 {
	return 9
}

func (*CallEvent) ValueSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) ValueInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.ValueSinceVersion()
}

func (*CallEvent) ValueDeprecated() uint16 {
	return 0
}

func (*CallEvent) ValueMetaAttribute(meta int) string {
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

func (*CallEvent) ValueMinValue() uint8 {
	return 0
}

func (*CallEvent) ValueMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*CallEvent) ValueNullValue() uint8 {
	return math.MaxUint8
}

func (*CallEvent) OutputSizeId() uint16 {
	return 10
}

func (*CallEvent) OutputSizeSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) OutputSizeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.OutputSizeSinceVersion()
}

func (*CallEvent) OutputSizeDeprecated() uint16 {
	return 0
}

func (*CallEvent) OutputSizeMetaAttribute(meta int) string {
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

func (*CallEvent) OutputSizeMinValue() uint32 {
	return 0
}

func (*CallEvent) OutputSizeMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*CallEvent) OutputSizeNullValue() uint32 {
	return math.MaxUint32
}

func (*CallEvent) RevertedId() uint16 {
	return 11
}

func (*CallEvent) RevertedSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) RevertedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.RevertedSinceVersion()
}

func (*CallEvent) RevertedDeprecated() uint16 {
	return 0
}

func (*CallEvent) RevertedMetaAttribute(meta int) string {
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

func (*CallEvent) RevertedMinValue() uint8 {
	return 0
}

func (*CallEvent) RevertedMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*CallEvent) RevertedNullValue() uint8 {
	return math.MaxUint8
}

func (*CallEvent) InputMetaAttribute(meta int) string {
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

func (*CallEvent) InputSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) InputInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.InputSinceVersion()
}

func (*CallEvent) InputDeprecated() uint16 {
	return 0
}

func (CallEvent) InputCharacterEncoding() string {
	return "null"
}

func (CallEvent) InputHeaderLength() uint64 {
	return 4
}

func (*CallEvent) OutputMetaAttribute(meta int) string {
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

func (*CallEvent) OutputSinceVersion() uint16 {
	return 0
}

func (c *CallEvent) OutputInActingVersion(actingVersion uint16) bool {
	return actingVersion >= c.OutputSinceVersion()
}

func (*CallEvent) OutputDeprecated() uint16 {
	return 0
}

func (CallEvent) OutputCharacterEncoding() string {
	return "null"
}

func (CallEvent) OutputHeaderLength() uint64 {
	return 4
}
