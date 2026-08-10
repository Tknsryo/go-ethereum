// Generated SBE (Simple Binary Encoding) message codec

package ethereum_tracing

import (
	"fmt"
	"io"
	"reflect"
)

type StatusEnum uint8
type StatusValues struct {
	Success   StatusEnum
	Error     StatusEnum
	NullValue StatusEnum
}

var Status = StatusValues{0, 1, 255}

func (s StatusEnum) Encode(_m *SbeGoMarshaller, _w io.Writer) error {
	if err := _m.WriteUint8(_w, uint8(s)); err != nil {
		return err
	}
	return nil
}

func (s *StatusEnum) Decode(_m *SbeGoMarshaller, _r io.Reader, actingVersion uint16) error {
	if err := _m.ReadUint8(_r, (*uint8)(s)); err != nil {
		return err
	}
	return nil
}

func (s StatusEnum) RangeCheck(actingVersion uint16, schemaVersion uint16) error {
	if actingVersion > schemaVersion {
		return nil
	}
	value := reflect.ValueOf(Status)
	for idx := 0; idx < value.NumField(); idx++ {
		if s == value.Field(idx).Interface() {
			return nil
		}
	}
	return fmt.Errorf("Range check failed on Status, unknown enumeration value %d", s)
}

func (*StatusEnum) EncodedLength() int64 {
	return 1
}

func (*StatusEnum) SuccessSinceVersion() uint16 {
	return 0
}

func (s *StatusEnum) SuccessInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.SuccessSinceVersion()
}

func (*StatusEnum) SuccessDeprecated() uint16 {
	return 0
}

func (*StatusEnum) ErrorSinceVersion() uint16 {
	return 0
}

func (s *StatusEnum) ErrorInActingVersion(actingVersion uint16) bool {
	return actingVersion >= s.ErrorSinceVersion()
}

func (*StatusEnum) ErrorDeprecated() uint16 {
	return 0
}
