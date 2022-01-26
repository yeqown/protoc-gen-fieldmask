package fieldmask

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// NestedFieldMask helps to build a nested field mask.
type NestedFieldMask map[string]NestedFieldMask

// New returns a new NestedFieldMask.
func New(fm *fieldmaskpb.FieldMask) NestedFieldMask {
	if fm == nil {
		return nil
	}

	paths := fm.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	nfm := make(NestedFieldMask, len(paths))
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if arr := strings.Split(p, "."); len(arr) > 1 {
			cur := nfm
			ok := false

			for _, p1 := range arr {
				if _, ok = cur[p1]; !ok {
					cur[p1] = make(NestedFieldMask, 1)
				}
				cur = cur[p1]
			}
		}

		if _, ok := nfm[p]; !ok {
			nfm[p] = make(NestedFieldMask, 1)
		}
	}

	return nfm
}

func (mask NestedFieldMask) Filter(m proto.Message) {
	if len(mask) == 0 {
		return
	}

	pr := m.ProtoReflect()
	pr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		nfm, ok := mask[string(fd.Name())]
		if !ok {
			pr.Clear(fd)
			return true
		}

		if len(nfm) == 0 {
			return true
		}

		// repeated field type
		if fd.IsList() {
			// TODO(@yeqown): skip basic type composite list fields
			l := pr.Get(fd).List()
			for i := 0; i < l.Len(); i++ {
				nfm.Filter(l.Get(i).Message().Interface())
			}
		} else if fd.Kind() == protoreflect.MessageKind {
			nfm.Filter(pr.Get(fd).Message().Interface())
		}

		return true
	})
}

func (mask NestedFieldMask) Prune(m proto.Message) {
	if len(mask) == 0 {
		return
	}

	pr := m.ProtoReflect()
	pr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		nfm, ok := mask[string(fd.Name())]
		if !ok {
			return true
		}

		if len(nfm) == 0 {
			pr.Clear(fd)
			return true
		}

		// repeated field type
		if fd.IsList() {
			// TODO(@yeqown): skip basic type composite list fields
			list := pr.Get(fd).List()
			for i := 0; i < list.Len(); i++ {
				nfm.Prune(list.Get(i).Message().Interface())
			}
		} else if fd.Kind() == protoreflect.MessageKind {
			nfm.Prune(pr.Get(fd).Message().Interface())
		}

		return true
	})
}
