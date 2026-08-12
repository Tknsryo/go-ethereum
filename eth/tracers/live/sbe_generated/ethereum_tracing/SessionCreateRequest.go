// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionCreateRequest struct {
	MessageType   MessageTypeEnum
	FilterMask    uint32
	AddressFilter []SessionCreateRequestAddressFilter
	TopicsFilter  []SessionCreateRequestTopicsFilter
}
type SessionCreateRequestAddressFilter struct {
	Address [20]uint8
}
type SessionCreateRequestTopicsFilter struct {
	Topics []SessionCreateRequestTopicsFilterTopics
}
type SessionCreateRequestTopicsFilterTopics struct {
	Topic [32]uint8
}

func (s *SessionCreateRequest) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := s.MessageType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, s.FilterMask); err != nil {
		return err
	}
	var AddressFilterBlockLength uint16 = 20
	if err := _m.WriteUint16(_w, AddressFilterBlockLength); err != nil {
		return err
	}
	var AddressFilterNumInGroup uint8 = uint8(len(s.AddressFilter))
	if err := _m.WriteUint8(_w, AddressFilterNumInGroup); err != nil {
		return err
	}
	for i := range s.AddressFilter {
		if err := s.AddressFilter[i].Encode(_m, _w); err != nil {
			return err
		}
	}
	var TopicsFilterBlockLength uint16 = 0
	if err := _m.WriteUint16(_w, TopicsFilterBlockLength); err != nil {
		return err
	}
	var TopicsFilterNumInGroup uint8 = uint8(len(s.TopicsFilter))
	if err := _m.WriteUint8(_w, TopicsFilterNumInGroup); err != nil {
		return err
	}
	for i := range s.TopicsFilter {
		if err := s.TopicsFilter[i].Encode(_m, _w); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionCreateRequest) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if s.MessageTypeInActingVersion(actingVersion) {
		if err := s.MessageType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !s.FilterMaskInActingVersion(actingVersion) {
		s.FilterMask = s.FilterMaskNullValue()
	} else {
		if err := _m.ReadUint32(_r, &s.FilterMask); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.AddressFilterInActingVersion(actingVersion) {
		var AddressFilterBlockLength uint16
		if err := _m.ReadUint16(_r, &AddressFilterBlockLength); err != nil {
			return err
		}
		var AddressFilterNumInGroup uint8
		if err := _m.ReadUint8(_r, &AddressFilterNumInGroup); err != nil {
			return err
		}
		if cap(s.AddressFilter) < int(AddressFilterNumInGroup) {
			s.AddressFilter = make([]SessionCreateRequestAddressFilter, AddressFilterNumInGroup)
		}
		s.AddressFilter = s.AddressFilter[:AddressFilterNumInGroup]
		for i := range s.AddressFilter {
			if err := s.AddressFilter[i].Decode(_m, _r, actingVersion, uint(AddressFilterBlockLength)); err != nil {
				return err
			}
		}
	}

	if s.TopicsFilterInActingVersion(actingVersion) {
		var TopicsFilterBlockLength uint16
		if err := _m.ReadUint16(_r, &TopicsFilterBlockLength); err != nil {
			return err
		}
		var TopicsFilterNumInGroup uint8
		if err := _m.ReadUint8(_r, &TopicsFilterNumInGroup); err != nil {
			return err
		}
		if cap(s.TopicsFilter) < int(TopicsFilterNumInGroup) {
			s.TopicsFilter = make([]SessionCreateRequestTopicsFilter, TopicsFilterNumInGroup)
		}
		s.TopicsFilter = s.TopicsFilter[:TopicsFilterNumInGroup]
		for i := range s.TopicsFilter {
			if err := s.TopicsFilter[i].Decode(_m, _r, actingVersion, uint(TopicsFilterBlockLength)); err != nil {
				return err
			}
		}
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionCreateRequest) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if err := s.MessageType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if s.FilterMaskInActingVersion(actingVersion) {
		if s.FilterMask < s.FilterMaskMinValue() || s.FilterMask > s.FilterMaskMaxValue() {
			return fmt.Errorf("Range check failed on s.FilterMask (%v < %v > %v)", s.FilterMaskMinValue(), s.FilterMask, s.FilterMaskMaxValue())
		}
	}
	for i := range s.AddressFilter {
		if err := s.AddressFilter[i].RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	for i := range s.TopicsFilter {
		if err := s.TopicsFilter[i].RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	return nil
}

func SessionCreateRequestInit(s *SessionCreateRequest) {
	return
}

func (s *SessionCreateRequestAddressFilter) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteBytes(_w, s.Address[:]); err != nil {
		return err
	}
	return nil
}

func (s *SessionCreateRequestAddressFilter) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !s.AddressInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			s.Address[idx] = s.AddressNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, s.Address[:]); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	return nil
}

func (s *SessionCreateRequestAddressFilter) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.AddressInActingVersion(actingVersion) {
		for idx := 0; idx < 20; idx++ {
			if s.Address[idx] < s.AddressMinValue() || s.Address[idx] > s.AddressMaxValue() {
				return fmt.Errorf("Range check failed on s.Address[%d] (%v < %v > %v)", idx, s.AddressMinValue(), s.Address[idx], s.AddressMaxValue())
			}
		}
	}
	return nil
}

func SessionCreateRequestAddressFilterInit(s *SessionCreateRequestAddressFilter) {
	return
}

func (s *SessionCreateRequestTopicsFilter) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	var TopicsBlockLength uint16 = 32
	if err := _m.WriteUint16(_w, TopicsBlockLength); err != nil {
		return err
	}
	var TopicsNumInGroup uint8 = uint8(len(s.Topics))
	if err := _m.WriteUint8(_w, TopicsNumInGroup); err != nil {
		return err
	}
	for i := range s.Topics {
		if err := s.Topics[i].Encode(_m, _w); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionCreateRequestTopicsFilter) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.TopicsInActingVersion(actingVersion) {
		var TopicsBlockLength uint16
		if err := _m.ReadUint16(_r, &TopicsBlockLength); err != nil {
			return err
		}
		var TopicsNumInGroup uint8
		if err := _m.ReadUint8(_r, &TopicsNumInGroup); err != nil {
			return err
		}
		if cap(s.Topics) < int(TopicsNumInGroup) {
			s.Topics = make([]SessionCreateRequestTopicsFilterTopics, TopicsNumInGroup)
		}
		s.Topics = s.Topics[:TopicsNumInGroup]
		for i := range s.Topics {
			if err := s.Topics[i].Decode(_m, _r, actingVersion, uint(TopicsBlockLength)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SessionCreateRequestTopicsFilter) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	for i := range s.Topics {
		if err := s.Topics[i].RangeCheck(actingVersion, schemaVersion); err != nil {
			return err
		}
	}
	return nil
}

func SessionCreateRequestTopicsFilterInit(s *SessionCreateRequestTopicsFilter) {
	return
}

func (s *SessionCreateRequestTopicsFilterTopics) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteBytes(_w, s.Topic[:]); err != nil {
		return err
	}
	return nil
}

func (s *SessionCreateRequestTopicsFilterTopics) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint) error {
	if !s.TopicInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			s.Topic[idx] = s.TopicNullValue()
		}
	} else {
		if err := _m.ReadBytes(_r, s.Topic[:]); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}
	return nil
}

func (s *SessionCreateRequestTopicsFilterTopics) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if s.TopicInActingVersion(actingVersion) {
		for idx := 0; idx < 32; idx++ {
			if s.Topic[idx] < s.TopicMinValue() || s.Topic[idx] > s.TopicMaxValue() {
				return fmt.Errorf("Range check failed on s.Topic[%d] (%v < %v > %v)", idx, s.TopicMinValue(), s.Topic[idx], s.TopicMaxValue())
			}
		}
	}
	return nil
}

func SessionCreateRequestTopicsFilterTopicsInit(s *SessionCreateRequestTopicsFilterTopics) {
	return
}

func (*SessionCreateRequest) SbeBlockLength() (blockLength uint16) {
	return 5
}

func (*SessionCreateRequest) SbeTemplateId() (templateId uint16) {
	return 20
}

func (*SessionCreateRequest) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*SessionCreateRequest) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*SessionCreateRequest) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionCreateRequest) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*SessionCreateRequest) MessageTypeId() uint16 {
	return 1
}

func (*SessionCreateRequest) MessageTypeSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequest) MessageTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.MessageTypeSinceVersion()
}

func (*SessionCreateRequest) MessageTypeDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequest) MessageTypeMetaAttribute(meta int) string {
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

func (*SessionCreateRequest) FilterMaskId() uint16 {
	return 2
}

func (*SessionCreateRequest) FilterMaskSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequest) FilterMaskInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.FilterMaskSinceVersion()
}

func (*SessionCreateRequest) FilterMaskDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequest) FilterMaskMetaAttribute(meta int) string {
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

func (*SessionCreateRequest) FilterMaskMinValue() uint32 {
	return 0
}

func (*SessionCreateRequest) FilterMaskMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*SessionCreateRequest) FilterMaskNullValue() uint32 {
	return math.MaxUint32
}

func (*SessionCreateRequestAddressFilter) AddressId() uint16 {
	return 1
}

func (*SessionCreateRequestAddressFilter) AddressSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequestAddressFilter) AddressInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.AddressSinceVersion()
}

func (*SessionCreateRequestAddressFilter) AddressDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequestAddressFilter) AddressMetaAttribute(meta int) string {
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

func (*SessionCreateRequestAddressFilter) AddressMinValue() uint8 {
	return 0
}

func (*SessionCreateRequestAddressFilter) AddressMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*SessionCreateRequestAddressFilter) AddressNullValue() uint8 {
	return math.MaxUint8
}

func (*SessionCreateRequestTopicsFilterTopics) TopicId() uint16 {
	return 1
}

func (*SessionCreateRequestTopicsFilterTopics) TopicSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequestTopicsFilterTopics) TopicInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TopicSinceVersion()
}

func (*SessionCreateRequestTopicsFilterTopics) TopicDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequestTopicsFilterTopics) TopicMetaAttribute(meta int) string {
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

func (*SessionCreateRequestTopicsFilterTopics) TopicMinValue() uint8 {
	return 0
}

func (*SessionCreateRequestTopicsFilterTopics) TopicMaxValue() uint8 {
	return math.MaxUint8 - 1
}

func (*SessionCreateRequestTopicsFilterTopics) TopicNullValue() uint8 {
	return math.MaxUint8
}

func (*SessionCreateRequest) AddressFilterId() uint16 {
	return 3
}

func (*SessionCreateRequest) AddressFilterSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequest) AddressFilterInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.AddressFilterSinceVersion()
}

func (*SessionCreateRequest) AddressFilterDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequestAddressFilter) SbeBlockLength() (blockLength uint) {
	return 20
}

func (*SessionCreateRequestAddressFilter) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*SessionCreateRequest) TopicsFilterId() uint16 {
	return 4
}

func (*SessionCreateRequest) TopicsFilterSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequest) TopicsFilterInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TopicsFilterSinceVersion()
}

func (*SessionCreateRequest) TopicsFilterDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequestTopicsFilter) SbeBlockLength() (blockLength uint) {
	return 0
}

func (*SessionCreateRequestTopicsFilter) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*SessionCreateRequestTopicsFilter) TopicsId() uint16 {
	return 5
}

func (*SessionCreateRequestTopicsFilter) TopicsSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequestTopicsFilter) TopicsInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.TopicsSinceVersion()
}

func (*SessionCreateRequestTopicsFilter) TopicsDeprecated() uint16 {
	return 0
}

func (*SessionCreateRequestTopicsFilterTopics) SbeBlockLength() (blockLength uint) {
	return 32
}

func (*SessionCreateRequestTopicsFilterTopics) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}
