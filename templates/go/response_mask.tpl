{{ define "response_mask" }}
{{ $inMessage := .InMessage }}
{{ $outMessage := .OutMessage }}
{{ $inMessageName := $inMessage.Name }}
{{ $outMessageName := $outMessage.Name }}
{{ $fmField := .FieldMaskField }}
{{ $method := .Method }}
{{ $methodName := $method.Name }}
{{ $outPkgPrefix := "" }}
{{ if $outMessagePkgName := $.OutMessagePkgName }}{{ $outPkgPrefix = printf "%s." $outMessagePkgName }}{{ end }}

// {{ $inMessageName }}_Res handles response field masking
type {{ $inMessageName }}_Res struct {
	fm *{{ $methodName }}_FieldMask
}

// Apply applies the field mask to the response message
func (r *{{ $inMessageName }}_Res) Apply(resp *{{ $outPkgPrefix }}{{ $outMessageName }}) {
	if r.fm.mode == pgfmpb.MaskMode_PRUNE {
		pgfmpb.Prune(resp, r.fm.res)
		return
	}
	pgfmpb.Filter(resp, r.fm.res)
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

// Mask{{ $fieldName }}_{{ $nestedFieldName }} marks the {{ $nestedFieldName }} field in nested {{ $fieldName }}
func (r *{{ $inMessageName }}_Res) Mask{{ $fieldName }}_{{ $nestedFieldName }}() {
	r.fm.setResFieldPath("{{ $field.Name }}.{{ $nestedField.Name }}")
}

// Masked{{ $fieldName }}_{{ $nestedFieldName }} checks if {{ $nestedFieldName }} is marked in the response
func (r *{{ $inMessageName }}_Res) Masked{{ $fieldName }}_{{ $nestedFieldName }}() bool {
	return r.fm.hasResFieldPath("{{ $field.Name }}.{{ $nestedField.Name }}")
}

				{{ end }}
			{{ end }}
		{{ else }}

// Mask{{ $fieldName }} marks the {{ $fieldName }} field in the response
func (r *{{ $inMessageName }}_Res) Mask{{ $fieldName }}() {
	r.fm.setResFieldPath("{{ $field.Name }}")
}

// Masked{{ $fieldName }} checks if {{ $fieldName }} is marked in the response
func (r *{{ $inMessageName }}_Res) Masked{{ $fieldName }}() bool {
	return r.fm.hasResFieldPath("{{ $field.Name }}")
}

		{{ end }}
	{{ end }}
{{ end }}

{{ end }}
