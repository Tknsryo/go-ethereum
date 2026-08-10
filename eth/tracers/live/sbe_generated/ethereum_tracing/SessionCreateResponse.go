// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionCreateResponse struct {
	MessageType MessageTypeEnum
	SessionID   uint32
	Status      StatusEnum
	PubEndpoint []uint8
	ErrorMsg    []uint8
}

func (s *SessionCreateResponse) Encode(_m *SbeGoMarshaller, _w io.Writer, doRangeCheck bool) error {
	if doRangeCheck {
		if err := s.RangeCheck(s.SbeSchemaVersion(), s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	if err := s.MessageType.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, s.SessionID); err != nil {
		return err
	}
	if err := s.Status.Encode(_m, _w); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.PubEndpoint))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.PubEndpoint); err != nil {
		return err
	}
	if err := _m.WriteUint32(_w, uint32(len(s.ErrorMsg))); err != nil {
		return err
	}
	if err := _m.WriteBytes(_w, s.ErrorMsg); err != nil {
		return err
	}
	return nil
}

func (s *SessionCreateResponse) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16, blockLength uint16, doRangeCheck bool) error {
	if s.MessageTypeInActingVersion(actingVersion) {
		if err := s.MessageType.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if !s.SessionIDInActingVersion(actingVersion) {
		s.SessionID = s.SessionIDNullValue()
	} else {
		if err := _m.ReadUint32(_r, &s.SessionID); err != nil {
			return err
		}
	}
	if s.StatusInActingVersion(actingVersion) {
		if err := s.Status.Decode(_m, _r, actingVersion); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
	}

	if s.PubEndpointInActingVersion(actingVersion) {
		var PubEndpointLength uint32
		if err := _m.ReadUint32(_r, &PubEndpointLength); err != nil {
			return err
		}
		if cap(s.PubEndpoint) < int(PubEndpointLength) {
			s.PubEndpoint = make([]uint8, PubEndpointLength)
		}
		s.PubEndpoint = s.PubEndpoint[:PubEndpointLength]
		if err := _m.ReadBytes(_r, s.PubEndpoint); err != nil {
			return err
		}
	}

	if s.ErrorMsgInActingVersion(actingVersion) {
		var ErrorMsgLength uint32
		if err := _m.ReadUint32(_r, &ErrorMsgLength); err != nil {
			return err
		}
		if cap(s.ErrorMsg) < int(ErrorMsgLength) {
			s.ErrorMsg = make([]uint8, ErrorMsgLength)
		}
		s.ErrorMsg = s.ErrorMsg[:ErrorMsgLength]
		if err := _m.ReadBytes(_r, s.ErrorMsg); err != nil {
			return err
		}
	}
	if doRangeCheck {
		if err := s.RangeCheck(actingVersion, s.SbeSchemaVersion()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SessionCreateResponse) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if err := s.MessageType.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	if s.SessionIDInActingVersion(actingVersion) {
		if s.SessionID < s.SessionIDMinValue() || s.SessionID > s.SessionIDMaxValue() {
			return fmt.Errorf("Range check failed on s.SessionID (%v < %v > %v)", s.SessionIDMinValue(), s.SessionID, s.SessionIDMaxValue())
		}
	}
	if err := s.Status.RangeCheck(actingVersion, schemaVersion); err != nil {
		return err
	}
	return nil
}

func SessionCreateResponseInit(s *SessionCreateResponse) {
	return
}

func (*SessionCreateResponse) SbeBlockLength() (blockLength uint16) {
	return 6
}

func (*SessionCreateResponse) SbeTemplateId() (templateId uint16) {
	return 21
}

func (*SessionCreateResponse) SbeSchemaId() (schemaId uint16) {
	return 1
}

func (*SessionCreateResponse) SbeSchemaVersion() (schemaVersion uint16) {
	return 1
}

func (*SessionCreateResponse) SbeSemanticType() (semanticType []byte) {
	return []byte("")
}

func (*SessionCreateResponse) SbeSemanticVersion() (semanticVersion string) {
	return "1.0.0"
}

func (*SessionCreateResponse) MessageTypeId() uint16 {
	return 1
}

func (*SessionCreateResponse) MessageTypeSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateResponse) MessageTypeInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.MessageTypeSinceVersion()
}

func (*SessionCreateResponse) MessageTypeDeprecated() uint16 {
	return 0
}

func (*SessionCreateResponse) MessageTypeMetaAttribute(meta int) string {
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

func (*SessionCreateResponse) SessionIDId() uint16 {
	return 2
}

func (*SessionCreateResponse) SessionIDSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateResponse) SessionIDInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SessionIDSinceVersion()
}

func (*SessionCreateResponse) SessionIDDeprecated() uint16 {
	return 0
}

func (*SessionCreateResponse) SessionIDMetaAttribute(meta int) string {
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

func (*SessionCreateResponse) SessionIDMinValue() uint32 {
	return 0
}

func (*SessionCreateResponse) SessionIDMaxValue() uint32 {
	return math.MaxUint32 - 1
}

func (*SessionCreateResponse) SessionIDNullValue() uint32 {
	return math.MaxUint32
}

func (*SessionCreateResponse) StatusId() uint16 {
	return 3
}

func (*SessionCreateResponse) StatusSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateResponse) StatusInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.StatusSinceVersion()
}

func (*SessionCreateResponse) StatusDeprecated() uint16 {
	return 0
}

func (*SessionCreateResponse) StatusMetaAttribute(meta int) string {
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

func (*SessionCreateResponse) PubEndpointMetaAttribute(meta int) string {
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

func (*SessionCreateResponse) PubEndpointSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateResponse) PubEndpointInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.PubEndpointSinceVersion()
}

func (*SessionCreateResponse) PubEndpointDeprecated() uint16 {
	return 0
}

func (SessionCreateResponse) PubEndpointCharacterEncoding() string {
	return "null"
}

func (SessionCreateResponse) PubEndpointHeaderLength() uint64 {
	return 4
}

func (*SessionCreateResponse) ErrorMsgMetaAttribute(meta int) string {
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

func (*SessionCreateResponse) ErrorMsgSinceVersion() uint16 {
	return 0
}

func (s *SessionCreateResponse) ErrorMsgInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ErrorMsgSinceVersion()
}

func (*SessionCreateResponse) ErrorMsgDeprecated() uint16 {
	return 0
}

func (SessionCreateResponse) ErrorMsgCharacterEncoding() string {
	return "null"
}

func (SessionCreateResponse) ErrorMsgHeaderLength() uint64 {
	return 4
}
