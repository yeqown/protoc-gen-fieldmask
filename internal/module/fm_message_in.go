package module

import (
	pbfieldmask "github.com/yeqown/protoc-gen-fieldmask/proto/fieldmask"

	pgs "github.com/lyft/protoc-gen-star"
)

type checkInMessageVO struct {
	FieldMaskOption *pbfieldmask.FieldMask
	FieldMaskField  pgs.Field
}

func (r *checkInMessageVO) invalid() bool {
	return r == nil || r.FieldMaskOption == nil || r.FieldMaskField == nil
}

const (
	// google_protobuf_FieldMask represents a google.protobuf.FieldMask type.
	// FIXME(@yeqown): FieldMaskModule type name would not be a constant string.
	google_protobuf_FieldMask = ".google.protobuf.FieldMask"
)

func (m *FieldMaskModule) checkInMessage(message pgs.Message) (r *checkInMessageVO) {
	if message == nil {
		return nil
	}

	r = new(checkInMessageVO)
	fields := message.Fields()
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		if f.Type().ProtoType() != pgs.MessageT {
			// not message type, fast fail.
			continue
		}
		m.Debugf("field (%s.%s:%s) is checking deeper", m.Name(), f.Name(), f.Descriptor().GetTypeName())

		if f.Descriptor().GetTypeName() != google_protobuf_FieldMask {
			// field's type not match google.protobuf.FieldMask
			continue
		}

		// DONE(@yeqiang): parse and record fieldmask.option.Option
		opt := new(pbfieldmask.FieldMask)
		_, err := f.Extension(pbfieldmask.E_Option, &opt)
		if err != nil || opt == nil {
			return nil
		}

		m.Debugf("message (%s) hit", m.Name())
		r.FieldMaskOption = opt
		r.FieldMaskField = f

		return r
	}

	return nil
}
