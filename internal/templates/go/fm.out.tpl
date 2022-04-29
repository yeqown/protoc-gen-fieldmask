{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $outMessagePkgName := .OutMessagePkgName }}
{{ $fmFieldName := .FieldMaskField.Name.UpperCamelCase }}

{{ $messageName := .OutMessage.Name.UpperCamelCase }}

{{ if .GenOutMessageVar }}
// _fm_{{ $messageName }} is built in variable for {{ $messageName }} to call FieldMask.Append
var _fm_{{ $messageName }} = new({{if $outMessagePkgName }}{{ $outMessagePkgName }}.{{end}}{{ $messageName }})
{{ end }}
{{ template "message" dict "Message" .OutMessage "inMessageName" $inMessageName "fmFieldName" $fmFieldName "inOut" "Out" "suffix" "" "pathSuffix" "" "messageName" $messageName }}

// Mask only affects the fields in the {{ $inMessageName }}.
func (x *{{ $inMessageName }}_FieldMask) Mask(m *{{if $outMessagePkgName }}{{ $outMessagePkgName }}.{{end}}{{ $outMessageName }}) *{{if $outMessagePkgName }}{{ $outMessagePkgName }}.{{end}}{{ $outMessageName }} {
   switch x.maskMode {
   case pbfieldmask.MaskMode_Filter:
        x.mask.Filter(m)
   case pbfieldmask.MaskMode_Prune:
        x.mask.Prune(m)
   }

   return m
}
