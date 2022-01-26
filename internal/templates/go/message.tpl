{{ $fields := .Message.Fields }}
{{ $fmFieldName := .fmFieldName }}
{{ $inMessageName := .inMessageName }}
{{ $inOut := .inOut }}
{{ $suffix := .suffix }}
{{ $pathSuffix := .pathSuffix }}
{{ $messageName := .messageName }}

{{ range $idx, $f := $fields }}
    {{ $fieldName := $f.Name.UpperCamelCase }}
    {{ $fieldPathKey := printf "%s%s" $pathSuffix $f.Name }}
    {{ if eq $fieldName $fmFieldName }}
    {{ else }}
         {{ $maskFuncName := printf "Mask%s_%s%s" $inOut $suffix $fieldName }}
         {{ $maskedFuncName := printf "Masked%s_%s%s" $inOut $suffix $fieldName }}

        // {{ $maskFuncName }} append {{ $messageName }}.{{ $fieldName }} into
        // {{ $inMessageName }}.{{ $fmFieldName }}.
        func (x *{{ $inMessageName }}) {{ $maskFuncName }}() *{{ $inMessageName }} {
              if x.{{ $fmFieldName }} == nil {
                  x.{{ $fmFieldName }} = new(fieldmaskpb.FieldMask)
              }
              x.{{ $fmFieldName }}.Append(_fm_{{ $messageName }}, "{{ $fieldPathKey }}")

              return x
        }

        // {{ $maskedFuncName }} indicates the field {{ $inMessageName }}.{{ $fieldName }}
        // exists in the {{ $inMessageName }}.{{ $fmFieldName }} or not.
        func (x *{{ $inMessageName }}_FieldMask) {{ $maskedFuncName }}() bool {
              if x.mask == nil {
                  return false
              }

              return x.mask.Masked("{{ $fieldPathKey }}")
        }

        {{ $recursive := (isMessage $f) }}
        {{ if eq $recursive true }}
            {{ $suffix2 := printf "%s%s_" $suffix $fieldName }}
            {{ $pathSuffix2 := printf "%s%s." $pathSuffix $f.Name }}
            {{ template "message" dict "Message" $f.Type.Embed "inMessageName" $inMessageName  "fmFieldName" $fmFieldName "inOut" "Out" "suffix" $suffix2 "pathSuffix" $pathSuffix2 "messageName" $messageName }}
        {{ end }}

    {{ end }}
{{ end }}
