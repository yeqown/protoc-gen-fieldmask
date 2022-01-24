{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $fmField := .FieldMaskField }}

{{ range $idx, $f := .InMessage.Fields }}
    {{ if eq $f.Name.UpperCamelCase $fmField.Name.UpperCamelCase }}
    {{ else }}
        // MaskIn_{{ $f.Name.UpperCamelCase }} indicates append {{ $outMessageName }}.{{ $f.Name.UpperCamelCase }} into
        // {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }}.
        func (x *{{ $inMessageName }}) MaskIn_{{ $f.Name.UpperCamelCase }}() *{{ $inMessageName }} {
              if x.{{ $fmField.Name.UpperCamelCase}} == nil {
                  x.{{ $fmField.Name.UpperCamelCase }} = new(fieldmaskpb.FieldMask)
              }
              x.{{ $fmField.Name.UpperCamelCase}}.Append(new({{ $inMessageName }}), "{{ $f.Name }}")

              return x
        }
    {{ end }}
{{ end }}

{{ range $idx, $f := .InMessage.Fields }}
    {{ if eq $f.Name.UpperCamelCase $fmField.Name.UpperCamelCase }}
    {{ else }}
        // Masked_{{ $f.Name.UpperCamelCase }} indicates the field {{ $inMessageName }}.{{ $f.Name.UpperCamelCase }}
        // exists in the {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }} or not.
        func (x *{{ $inMessageName }}_FieldMask) MaskedIn_{{ $f.Name.UpperCamelCase }}() bool {
              if x.maskMapping == nil {
                  return false
              }

              _, ok := x.maskMapping["{{ $f.Name }}"]
              return ok
        }
    {{ end }}
{{ end }}