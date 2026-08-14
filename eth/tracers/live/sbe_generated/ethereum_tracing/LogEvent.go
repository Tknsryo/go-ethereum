// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type LogEvent struct {
	EventType EventTypeEnum
	Timestamp uint64
	Index     uint32
	Address   [20]uint8
	Removed   uint8
	Topics    []LogEventTopics
	LogData   []uint8
}
type LogEventTopics struct {
	Topic [32]uint8
}

func (l *LogEvent) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := l.RangeCheck(l.SbeSchemaVersion(), l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := l.EventType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, l.Timestamp); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, l.Index); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, l.Address[:]); err != nil {
		return err
	}
	if err := _m.WriteUint8(_w, l.Removed); err != nil {
		return err
	}
	var TopicsBlockLength uint16 = 32
	if err := _m.WriteUint16(_w, TopicsBlockLength); err != nil {
		return err
	}
	var TopicsNumInGroup uint8 = uint8(len(l.Topics))
	if err := _m.WriteUint8(_w, TopicsNumInGroup); err != nil {
		return err
	}
	for i := range l.Topics {
		if err := l.Topics[i].Encode(_m, _w); err != nil {
			return err
		}
	}
	if err := _m.WriteUint32(_w, uint32(len(l.LogData))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, l.LogData); err != nil {
		return err
	}
	return nil
}

func (l *LogEvent) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if l.EventTypeInActingVersion(actingVersion) {
		if err := l.EventType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !l.TimestampInActingVersion(actingVersion) {
		l.Timestamp = l.TimestampNullValue()
	} else {
		if err := _m.ReadUint64(_r, &l.Timestamp); err != nil {
			return err
		}
	}
	if !l.IndexInActingVersion(actingVersion) {
		l.Index = l.IndexNullValue()
	} else {
		if err := _m.ReadUint32(_r, &l.Index); err != nil {
			return err
		}
	}
	if !l.AddressInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			l.Address[idx] = l.AddressNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, l.Address[:]); err != nil {
			return err
		}
	}
	if !l.RemovedInActingVersion(actingVersion) {
		l.Removed = l.RemovedNullValue()
	} else {
		if err := _m.ReadUint8(_r, &l.Removed); err != nil {
			return err
		}
	}
	if actingVersion > l.SbeSchemaVersion() && blockLength > l.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-l.SbeBlockLength()))
	}

	if l.TopicsInActingVersion(actingVersion) {
		var TopicsBlockLength uint16
		if err := _m.ReadUint16(_r, &TopicsBlockLength); err != nil {
			return err
		}
		var TopicsNumInGroup uint8
		if err := _m.ReadUint8(_r, &TopicsNumInGroup); err != nil {
			return err
		}
		if cap(l.Topics) < int(TopicsNumInGroup) {
			l.Topics = make([]LogEventTopics, TopicsNumInGroup)
		}
		l.Topics = l.Topics[:TopicsNumInGroup]
		for i := range l.Topics {
			if err := l.Topics[i].Decode(_m, _r, actingVersion, uint(TopicsBlockLength)); err != nil {
				return err
			}
		}
	}

	if l.LogDataInActingVersion(actingVersion) {
		var LogDataLength uint32
		if err := _m.ReadUint32(_r, &LogDataLength); err != nil {
			return err
		}
		if cap(l.LogData) < int(LogDataLength) {
			l.LogData = make([]uint8, LogDataLength)
		}
		l.LogData = l.LogData[:LogDataLength]
		if err := _m.ReadBytes(_r, l.LogData); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := l.RangeCheck(actingVersion, l.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (l *LogEvent) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if err := l.EventType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if l.TimestampInActingVersion(actingVersion) {
		if l.Timestamp < l.TimestampMinValue() || l.Timestamp > l.TimestampMaxValue() {
			return fmt.Errorf("Range check failed on l.Timestamp (%v < %v > %v)", l.TimestampMinValue(), l.Timestamp, l.TimestampMaxValue())
		}
	}
	if l.IndexInActingVersion(actingVersion) {
		if l.Index < l.IndexMinValue() || l.Index > l.IndexMaxValue() {
			return fmt.Errorf("Range check failed on l.Index (%v < %v > %v)", l.IndexMinValue(), l.Index, l.IndexMaxValue())
		}
	}
	if l.AddressInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			if l.Address[idx] < l.AddressMinValue() || l.Address[idx] > l.AddressMaxValue() {
				return fmt.Errorf("Range check failed on l.Address[%d] (%v < %v > %v)", idx, l.AddressMinValue(), l.Address[idx], l.AddressMaxValue())
			}
		}
	}
	if l.RemovedInActingVersion(actingVersion) {
		if l.Removed < l.RemovedMinValue() || l.Removed > l.RemovedMaxValue() {
			return fmt.Errorf("Range check failed on l.Removed (%v < %v > %v)", l.RemovedMinValue(), l.Removed, l.RemovedMaxValue())
		}
	}
	for i := range l.Topics {
		if err := l.Topics[i].RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	return nil
}

func LogEventInit(l *LogEvent) {
	return
}

func (l *LogEventTopics) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteBytes(_w, l.Topic[:]); err != nil {
		return err
	}
	return nil
}

func (l *LogEventTopics) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !l.TopicInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			l.Topic[idx] = l.TopicNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, l.Topic[:]); err != nil {
			return err
		}
	}
	if actingVersion > l.SbeSchemaVersion() && blockLength > l.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-l.SbeBlockLength()))
	}
	return nil
}

func (l *LogEventTopics) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if l.TopicInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if l.Topic[idx] < l.TopicMinValue() || l.Topic[idx] > l.TopicMaxValue() {
				return fmt.Errorf("Range check failed on l.Topic[%d] (%v < %v > %v)", idx, l.TopicMinValue(), l.Topic[idx], l.TopicMaxValue())
			}
		}
	}
	return nil
}

func LogEventTopicsInit(l *LogEventTopics) {
	return
}

func (*LogEvent) SbeBlockLength() (blockLength uint16) {
	return 34
}

func (*LogEvent) SbeTemplateId() (templateId uint16) {
	return 15
}

func (*LogEvent) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*LogEvent) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*LogEvent) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*LogEvent) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*LogEvent) EventTypeId() uint16 {
	return 1
}

func (*LogEvent) EventTypeSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) EventTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.EventTypeSinceVersion()
}

func (*LogEvent) EventTypeDeprecated() uint16 {
	return 0
}

func (*LogEvent) EventTypeMetaAttribute(meta int) string {
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

func (*LogEvent) TimestampId() uint16 {
	return 2
}

func (*LogEvent) TimestampSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) TimestampInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.TimestampSinceVersion()
}

func (*LogEvent) TimestampDeprecated() uint16 {
	return 0
}

func (*LogEvent) TimestampMetaAttribute(meta int) string {
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

func (*LogEvent) TimestampMinValue() uint64 {
	return 0
}

func (*LogEvent) TimestampMaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*LogEvent) TimestampNullValue() uint64 {
	return math.MaxUint64
}

func (*LogEvent) IndexId() uint16 {
	return 3
}

func (*LogEvent) IndexSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) IndexInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.IndexSinceVersion()
}

func (*LogEvent) IndexDeprecated() uint16 {
	return 0
}

func (*LogEvent) IndexMetaAttribute(meta int) string {
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

func (*LogEvent) IndexMinValue() uint32 {
	return 0
}

func (*LogEvent) IndexMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*LogEvent) IndexNullValue() uint32 {
	return math.MaxUint32
}

func (*LogEvent) AddressId() uint16 {
	return 4
}

func (*LogEvent) AddressSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) AddressInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.AddressSinceVersion()
}

func (*LogEvent) AddressDeprecated() uint16 {
	return 0
}

func (*LogEvent) AddressMetaAttribute(meta int) string {
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

func (*LogEvent) AddressMinValue() uint8 {
	return 0
}

func (*LogEvent) AddressMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*LogEvent) AddressNullValue() uint8 {
	return math.MaxUint8
}

func (*LogEvent) RemovedId() uint16 {
	return 5
}

func (*LogEvent) RemovedSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) RemovedInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.RemovedSinceVersion()
}

func (*LogEvent) RemovedDeprecated() uint16 {
	return 0
}

func (*LogEvent) RemovedMetaAttribute(meta int) string {
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

func (*LogEvent) RemovedMinValue() uint8 {
	return 0
}

func (*LogEvent) RemovedMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*LogEvent) RemovedNullValue() uint8 {
	return math.MaxUint8
}

func (*LogEventTopics) TopicId() uint16 {
	return 1
}

func (*LogEventTopics) TopicSinceVersion() uint16 {
	return 0
}

func (l *LogEventTopics) TopicInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.TopicSinceVersion()
}

func (*LogEventTopics) TopicDeprecated() uint16 {
	return 0
}

func (*LogEventTopics) TopicMetaAttribute(meta int) string {
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

func (*LogEventTopics) TopicMinValue() uint8 {
	return 0
}

func (*LogEventTopics) TopicMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*LogEventTopics) TopicNullValue() uint8 {
	return math.MaxUint8
}

func (*LogEvent) TopicsId() uint16 {
	return 6
}

func (*LogEvent) TopicsSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) TopicsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.TopicsSinceVersion()
}

func (*LogEvent) TopicsDeprecated() uint16 {
	return 0
}

func (*LogEventTopics) SbeBlockLength() (blockLength uint) {
	return 32
}

func (*LogEventTopics) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*LogEvent) LogDataMetaAttribute(meta int) string {
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

func (*LogEvent) LogDataSinceVersion() uint16 {
	return 0
}

func (l *LogEvent) LogDataInActingVersion(actingVersion uint16) bool {
	return actingVersion >= l.LogDataSinceVersion()
}

func (*LogEvent) LogDataDeprecated() uint16 {
	return 0
}

func (LogEvent) LogDataCharacterEncoding() string {
	return "null"
}

func (LogEvent) LogDataHeaderLength() uint64 {
	return 4
}
