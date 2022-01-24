{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $fmField := .FieldMaskField }}

{{ range $idx, $f := .OutMessage.Fields }}
    // Masked_{{ $f.Name.UpperCamelCase }} indicates the field {{ $inMessageName }}.{{ $f.Name.UpperCamelCase }}
    // exists in the {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }} or not.
    func (x *{{ $inMessageName }}_FieldMask) MaskedOut_{{ $f.Name.UpperCamelCase }}() bool {
          if x.maskMapping == nil {
              return false
          }

          _, ok := x.maskMapping["{{ $f.Name }}"]
          return ok
    }
{{ end }}

// Mask only affects the fields in the {{ $inMessageName }}.
func (x *{{ $inMessageName }}_FieldMask) Mask(m *{{ $outMessageName }}) *{{ $outMessageName }} {
   switch x.maskMode {
   case pbfieldmask.MaskMode_Filter:
        x.filter(m)
   case pbfieldmask.MaskMode_Prune:
        x.prune(m)
   }

   return m
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