{{ $inMessageName := .InMessage.Name }}
{{ $fmFieldName := .FieldMaskField.Name.UpperCamelCase }}

{{ $messageName := .InMessage.Name.UpperCamelCase }}

// _fm_{{ $messageName }} is built in variable for {{ $messageName }} to call FieldMask.Append
var _fm_{{ $messageName }} = new({{ $messageName }})

{{ template "message" dict "Message" .InMessage "inMessageName" $inMessageName "fmFieldName" $fmFieldName "inOut" "In" "suffix" "" "pathSuffix" "" "messageName" $messageName }}