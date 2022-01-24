{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $fmField := .FieldMaskField }}

{{ range $idx, $f := .OutMessage.Fields }}
    // MaskOut_{{ $f.Name.UpperCamelCase }} indicates append {{ $outMessageName }}.{{ $f.Name.UpperCamelCase }} into
    // {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }}.
    func (x *{{ $inMessageName }}) MaskOut_{{ $f.Name.UpperCamelCase }}() *{{ $inMessageName }} {
          if x.{{ $fmField.Name.UpperCamelCase}} == nil {
              x.{{ $fmField.Name.UpperCamelCase }} = new(fieldmaskpb.FieldMask)
          }
          x.{{ $fmField.Name.UpperCamelCase}}.Append(new({{ $outMessageName }}), "{{ $f.Name }}")

          return x
    }
{{ end}}
