{{ define "request_mask" }}
{{ $inMessage := .InMessage }}
{{ $outMessage := .OutMessage }}
{{ $inMessageName := $inMessage.Name }}
{{ $outMessageName := $outMessage.Name }}
{{ $fmField := .FieldMaskField }}

// {{ $inMessageName }}_RequestMask handles request field masking
type {{ $inMessageName }}_RequestMask struct {
	fm *{{ $inMessageName }}
}

// Apply applies the field mask to the request message
func (rm *{{ $inMessageName }}_RequestMask) Apply(req *{{ $inMessageName }}) {
	// Implementation would depend on the specific field mask library being used
	// For now, we'll leave this as a placeholder
}

{{ range $idx, $field := $inMessage.Fields }}
	{{ if ne $field.Name.String $fmField.Name.String }}
		{{ $fieldOpts := fieldOptions $field }}
		{{ if not $fieldOpts.Ignore }}
			{{ $fieldName := $field.Name.UpperCamelCase }}

// {{ $fieldName }} marks the {{ $fieldName }} field in the request
func (rm *{{ $inMessageName }}_RequestMask) {{ $fieldName }}() *{{ $inMessageName }}_RequestMask {
	rm.fm.{{ $fmField.Name.UpperCamelCase }}.Paths = append(rm.fm.{{ $fmField.Name.UpperCamelCase }}.Paths, "{{ $field.Name }}")
	return rm
}

		{{ end }}
	{{ end }}
{{ end }}

{{ end }}