// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

type SessionCreateRequest struct {
	MessageType MessageTypeEnum
	FilterMask  uint32
	ClientID0   uint64
	ClientID1   uint64
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
	if err := _m.WriteUint64(_w, s.ClientID0); err != nil {
		return err
	}
	if err := _m.WriteUint64(_w, s.ClientID1); err != nil {
		return err
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
	if !s.ClientID0InActingVersion(actingVersion) {
		s.ClientID0 = s.ClientID0NullValue()
	} else {
		if err := _m.ReadUint64(_r, &s.ClientID0); err != nil {
			return err
		}
	}
	if !s.ClientID1InActingVersion(actingVersion) {
		s.ClientID1 = s.ClientID1NullValue()
	} else {
		if err := _m.ReadUint64(_r, &s.ClientID1); err != nil {
			return err
		}
	}
	if actingVersion > s.SbeSchemaVersion() && blockLength > s.SbeBlockLength() {
		io.CopyN(ioutil.Discard, _r, int64(blockLength-s.SbeBlockLength()))
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
	if s.ClientID0InActingVersion(actingVersion) {
		if s.ClientID0 < s.ClientID0MinValue() || s.ClientID0 > s.ClientID0MaxValue() {
			return fmt.Errorf("Range check failed on s.ClientID0 (%v < %v > %v)", s.ClientID0MinValue(), s.ClientID0, s.ClientID0MaxValue())
		}
	}
	if s.ClientID1InActingVersion(actingVersion) {
		if s.ClientID1 < s.ClientID1MinValue() || s.ClientID1 > s.ClientID1MaxValue() {
			return fmt.Errorf("Range check failed on s.ClientID1 (%v < %v > %v)", s.ClientID1MinValue(), s.ClientID1, s.ClientID1MaxValue())
		}
	}
	return nil
}

func SessionCreateRequestInit(s *SessionCreateRequest) {
	return
}

func (*SessionCreateRequest) SbeBlockLength() (blockLength uint16) {
	return 21
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

func (*SessionCreateRequest) ClientID0Id() uint16 {
	return 3
}

func (*SessionCreateRequest) ClientID0SinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequest) ClientID0InActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClientID0SinceVersion()
}

func (*SessionCreateRequest) ClientID0Deprecated() uint16 {
	return 0
}

func (*SessionCreateRequest) ClientID0MetaAttribute(meta int) string {
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

func (*SessionCreateRequest) ClientID0MinValue() uint64 {
	return 0
}

func (*SessionCreateRequest) ClientID0MaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*SessionCreateRequest) ClientID0NullValue() uint64 {
	return math.MaxUint64
}

func (*SessionCreateRequest) ClientID1Id() uint16 {
	return 4
}

func (*SessionCreateRequest) ClientID1SinceVersion() uint16 {
	return 0
}

func (s *SessionCreateRequest) ClientID1InActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ClientID1SinceVersion()
}

func (*SessionCreateRequest) ClientID1Deprecated() uint16 {
	return 0
}

func (*SessionCreateRequest) ClientID1MetaAttribute(meta int) string {
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

func (*SessionCreateRequest) ClientID1MinValue() uint64 {
	return 0
}

func (*SessionCreateRequest) ClientID1MaxValue() uint64 {
	return math.MaxUint64 - 1
}

func (*SessionCreateRequest) ClientID1NullValue() uint64 {
	return math.MaxUint64
}
