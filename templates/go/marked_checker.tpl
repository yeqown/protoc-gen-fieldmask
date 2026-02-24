{{ define "marked_checker" }}
{{ $inMessage := .InMessage }}
{{ $outMessage := .OutMessage }}
{{ $inMessageName := $inMessage.Name }}
{{ $outMessageName := $outMessage.Name }}
{{ $fmField := .FieldMaskField }}

// {{ $inMessageName }}_Masked provides field marking checks
type {{ $inMessageName }}_Masked struct {
	fm *{{ $inMessageName }}_FieldMask
}

// Request returns the request field marking checker
func (m *{{ $inMessageName }}_Masked) Request() *{{ $inMessageName }}_Masked_Req {
	return &{{ $inMessageName }}_Masked_Req{
		fm: m.fm,
	}
}

// Response returns the response field marking checker
func (m *{{ $inMessageName }}_Masked) Response() *{{ $inMessageName }}_Masked_Res {
	return &{{ $inMessageName }}_Masked_Res{
		fm: m.fm,
	}
}

// {{ $inMessageName }}_Masked_Req provides request field marking checks
type {{ $inMessageName }}_Masked_Req struct {
	fm *{{ $inMessageName }}_FieldMask
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

// Masked{{ $fieldName }}_{{ $nestedFieldName }} checks if the {{ $nestedFieldName }} field is marked in the request
func (mr *{{ $inMessageName }}_Masked_Req) Masked{{ $fieldName }}_{{ $nestedFieldName }}() bool {
	return mr.fm.hasReqFieldPath("{{ $field.Name }}.{{ $nestedField.Name }}")
}

					{{ end }}
				{{ end }}
			{{ else }}

// Masked{{ $fieldName }} checks if the {{ $fieldName }} field is marked in the request
func (mr *{{ $inMessageName }}_Masked_Req) Masked{{ $fieldName }}() bool {
	return mr.fm.hasReqFieldPath("{{ $field.Name }}")
}

			{{ end }}
		{{ end }}
	{{ end }}
{{ end }}

// {{ $inMessageName }}_Masked_Res provides response field marking checks
type {{ $inMessageName }}_Masked_Res struct {
	fm *{{ $inMessageName }}_FieldMask
}

{{ range $idx, $field := $outMessage.Fields }}
	{{ $fieldOpts := fieldOptions $field }}
	{{ if $fieldOpts.Mask }}
		{{ $fieldName := $field.Name.UpperCamelCase }}
		{{ $isNested := and $fieldOpts.Nested (isMessage $field) }}

		{{ if $isNested }}
			{{ range $nestedField := $field.Message.Fields }}
				{{ $nestedFieldOpts := fieldOptions $nestedField }}
				{{ if $nestedFieldOpts.Mask }}
					{{ $nestedFieldName := $nestedField.Name.UpperCamelCase }}

// Masked{{ $fieldName }}_{{ $nestedFieldName }} checks if the {{ $nestedFieldName }} field is marked in the response
func (mr *{{ $inMessageName }}_Masked_Res) Masked{{ $fieldName }}_{{ $nestedFieldName }}() bool {
	return mr.fm.hasResFieldPath("{{ $field.Name }}.{{ $nestedField.Name }}")
}

				{{ end }}
			{{ end }}
		{{ else }}

// Masked{{ $fieldName }} checks if the {{ $fieldName }} field is marked in the response
func (mr *{{ $inMessageName }}_Masked_Res) Masked{{ $fieldName }}() bool {
	return mr.fm.hasResFieldPath("{{ $field.Name }}")
}

		{{ end }}
	{{ end }}
{{ end }}

{{ end }}
