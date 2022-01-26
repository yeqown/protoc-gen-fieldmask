{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $fmField := .FieldMaskField }}

// _fm_{{ $outMessageName }} is built in variable for {{ $outMessageName }} to call FieldMask.Append
var _fm_{{ $outMessageName }} = new({{ $outMessageName }})

{{ range $idx, $f := .OutMessage.Fields }}
// MaskOut_{{ $f.Name.UpperCamelCase }} indicates append {{ $outMessageName }}.{{ $f.Name.UpperCamelCase }} into
// {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }}.
func (x *{{ $inMessageName }}) MaskOut_{{ $f.Name.UpperCamelCase }}() *{{ $inMessageName }} {
      if x.{{ $fmField.Name.UpperCamelCase}} == nil {
          x.{{ $fmField.Name.UpperCamelCase }} = new(fieldmaskpb.FieldMask)
      }
      x.{{ $fmField.Name.UpperCamelCase}}.Append(_fm_{{ $outMessageName }}, "{{ $f.Name }}")

      return x
}
{{ end}}

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
