{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $fmField := .FieldMaskField }}

func (x *{{ $inMessageName }}) FieldMaskWithMode(mode pbfieldmask.MaskMode) *{{ $inMessageName }}_FieldMask {
    fm := &{{ $inMessageName }}_FieldMask{
        maskMode: mode,
        maskMapping: make(map[string]struct{}, len(x.{{ $fmField.Name.UpperCamelCase }}.GetPaths())),
    }

    for _, path := range x.{{ $fmField.Name.UpperCamelCase }}.GetPaths() {
        fm.maskMapping[path] = struct{}{}
    }

    return fm
}

// FieldMask_Prune generates *{{ $inMessageName }}_FieldMask with filter mode, so that
// only the fields in the {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }} will be
// appended into the {{ $inMessageName }}.
func (x *{{ $inMessageName }}) FieldMask_Filter() *{{ $inMessageName }}_FieldMask {
	return x.FieldMaskWithMode(pbfieldmask.MaskMode_Filter)
}

// FieldMask_Prune generates *{{ $inMessageName }}_FieldMask with prune mode, so that
// only the fields NOT in the {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }} will be
// appended into the {{ $inMessageName }}.
func (x *{{ $inMessageName }}) FieldMask_Prune() *{{ $inMessageName }}_FieldMask {
	return x.FieldMaskWithMode(pbfieldmask.MaskMode_Prune)
}

// {{ $inMessageName }}_FieldMask provide provide helper functions to deal with FieldMask.
type {{ $inMessageName }}_FieldMask struct {
    maskMode pbfieldmask.MaskMode
    maskMapping map[string]struct{}
}

// filter will retain the fields those are in the maskMapping
func (x *{{ $inMessageName }}_FieldMask) filter(m proto.Message) {
    if len(x.maskMapping) == 0 {
        return
    }

    pr := m.ProtoReflect()
    pr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
        _, ok := x.maskMapping[string(fd.Name())]
        if !ok {
            pr.Clear(fd)
            return true
        }

        // TODO(@yeqown): support deeper fields masking
        return true
    })
}

// prune will remove fields those are in the maskMapping
func (x *{{ $inMessageName }}_FieldMask) prune(m proto.Message) {
    if len(x.maskMapping) == 0 {
        return
    }

    pr := m.ProtoReflect()
    pr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
        _, ok := x.maskMapping[string(fd.Name())]
        if !ok {
            return true
        }

        // TODO(@yeqown): support deeper fields masking
        pr.Clear(fd)
        return true
    })
}