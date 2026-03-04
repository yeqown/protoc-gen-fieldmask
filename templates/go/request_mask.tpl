{{ define "request_mask" }}
{{ $inMessage := .InMessage }}
{{ $outMessage := .OutMessage }}
{{ $inMessageName := $inMessage.Name }}
{{ $outMessageName := $outMessage.Name }}
{{ $fmField := .FieldMaskField }}
{{ $method := .Method }}
{{ $methodName := $method.Name }}

// {{ $inMessageName }}_Req handles request field masking
type {{ $inMessageName }}_Req struct {
	fm *{{ $methodName }}_FieldMask
}

{{ range $idx, $field := $inMessage.Fields }}
	{{ if ne $field.Name.String $fmField.Name.String }}
		{{ $fieldOpts := fieldOptions $field }}
		{{ if $fieldOpts.Mask }}
			{{ $fieldName := $field.Name.UpperCamelCase }}
			{{ $isNested := and $fieldOpts.Nested (isMessage $field) }}

			{{ if $isNested }}
				{{ range $nestedField := $field.Message.Fields }}
					{{ $nestedFieldOpts := fieldOptions $nestedField }}
					{{ if $nestedFieldOpts.Mask }}
						{{ $nestedFieldName := $nestedField.Name.UpperCamelCase }}

// Mask{{ $fieldName }}_{{ $nestedFieldName }} marks the {{ $nestedFieldName }} field in nested {{ $fieldName }}
func (r *{{ $inMessageName }}_Req) Mask{{ $fieldName }}_{{ $nestedFieldName }}() {
	r.fm.setReqFieldPath("{{ $field.Name }}.{{ $nestedField.Name }}")
}

// Masked{{ $fieldName }}_{{ $nestedFieldName }} checks if {{ $nestedFieldName }} is marked in the request
func (r *{{ $inMessageName }}_Req) Masked{{ $fieldName }}_{{ $nestedFieldName }}() bool {
	return r.fm.hasReqFieldPath("{{ $field.Name }}.{{ $nestedField.Name }}")
}

					{{ end }}
				{{ end }}
			{{ else }}

// Mask{{ $fieldName }} marks the {{ $fieldName }} field in the request
func (r *{{ $inMessageName }}_Req) Mask{{ $fieldName }}() {
	r.fm.setReqFieldPath("{{ $field.Name }}")
}

// Masked{{ $fieldName }} checks if {{ $fieldName }} is marked in the request
func (r *{{ $inMessageName }}_Req) Masked{{ $fieldName }}() bool {
	return r.fm.hasReqFieldPath("{{ $field.Name }}")
}

			{{ end }}
		{{ end }}
	{{ end }}
{{ end }}

{{ end }}
