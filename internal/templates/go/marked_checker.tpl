{{ define "marked_checker" }}
{{ $inMessage := .InMessage }}
{{ $outMessage := .OutMessage }}
{{ $inMessageName := $inMessage.Name }}
{{ $outMessageName := $outMessage.Name }}
{{ $fmField := .FieldMaskField }}

// {{ $inMessageName }}_Marked provides field marking checks
type {{ $inMessageName }}_Marked struct {
	fm *{{ $inMessageName }}
}

// Response returns the response field marking checker
func (m *{{ $inMessageName }}_Marked) Response() *{{ $inMessageName }}_Marked_Response {
	return &{{ $inMessageName }}_Marked_Response{
		marked: m,
	}
}

{{ range $idx, $field := $inMessage.Fields }}
	{{ if ne $field.Name.String $fmField.Name.String }}
		{{ $fieldOpts := fieldOptions $field }}
		{{ if not $fieldOpts.Ignore }}
			{{ $fieldName := $field.Name.UpperCamelCase }}

// {{ $fieldName }} checks if the {{ $fieldName }} field is marked in the request
func (m *{{ $inMessageName }}_Marked) {{ $fieldName }}() bool {
	for _, path := range m.fm.{{ $fmField.Name.UpperCamelCase }}.Paths {
		if path == "{{ $field.Name }}" {
			return true
		}
	}
	return false
}

		{{ end }}
	{{ end }}
{{ end }}

// {{ $inMessageName }}_Marked_Response provides response field marking checks
type {{ $inMessageName }}_Marked_Response struct {
	marked *{{ $inMessageName }}_Marked
}

{{ range $idx, $field := $outMessage.Fields }}
	{{ $fieldOpts := fieldOptions $field }}
	{{ if not $fieldOpts.Ignore }}
		{{ $fieldName := $field.Name.UpperCamelCase }}

// {{ $fieldName }} checks if the {{ $fieldName }} field is marked in the response
func (mr *{{ $inMessageName }}_Marked_Response) {{ $fieldName }}() bool {
	for _, path := range mr.marked.fm.{{ $fmField.Name.UpperCamelCase }}.Paths {
		if path == "{{ $field.Name }}" {
			return true
		}
	}
	return false
}

	{{ end }}
{{ end }}

{{ end }}