{{ define "response_mask" }}
{{ $inMessage := .InMessage }}
{{ $outMessage := .OutMessage }}
{{ $inMessageName := $inMessage.Name }}
{{ $outMessageName := $outMessage.Name }}
{{ $fmField := .FieldMaskField }}
{{ $outPkgPrefix := "" }}
{{ if $outMessagePkgName := $.OutMessagePkgName }}{{ $outPkgPrefix = printf "%s." $outMessagePkgName }}{{ end }}

// {{ $inMessageName }}_ResponseMask handles response field masking
type {{ $inMessageName }}_ResponseMask struct {
	fm                *{{ $inMessageName }}
	outMessageName    string
	outMessagePkgName string
}

// Apply applies the field mask to the response message
func (rm *{{ $inMessageName }}_ResponseMask) Apply(resp *{{ $outPkgPrefix }}{{ $outMessageName }}) {
	// Implementation would depend on the specific field mask library being used
	// For now, we'll leave this as a placeholder
}

{{ range $idx, $field := $outMessage.Fields }}
	{{ $fieldOpts := fieldOptions $field }}
	{{ if not $fieldOpts.Ignore }}
		{{ $fieldName := $field.Name.UpperCamelCase }}

// {{ $fieldName }} marks the {{ $fieldName }} field in the response
func (rm *{{ $inMessageName }}_ResponseMask) {{ $fieldName }}() *{{ $inMessageName }}_ResponseMask {
	path := "{{ $field.Name }}"
	rm.fm.{{ $fmField.Name.UpperCamelCase }}.Paths = append(rm.fm.{{ $fmField.Name.UpperCamelCase }}.Paths, path)
	return rm
}

		{{ if and $fieldOpts.Nested (isMessage $field) }}
// {{ $fieldName }} returns a nested field mask for {{ $fieldName }}
func (rm *{{ $inMessageName }}_ResponseMask) {{ $fieldName }}() *{{ $inMessageName }}_ResponseMask_{{ $fieldName }} {
	return &{{ $inMessageName }}_ResponseMask_{{ $fieldName }}{
		parent: rm,
		path:   "{{ $field.Name }}",
	}
}

// {{ $inMessageName }}_ResponseMask_{{ $fieldName }} handles nested field masking for {{ $fieldName }}
type {{ $inMessageName }}_ResponseMask_{{ $fieldName }} struct {
	parent *{{ $inMessageName }}_ResponseMask
	path   string
}

// FieldMask returns a field mask for the nested {{ $fieldName }} message
func (nm *{{ $inMessageName }}_ResponseMask_{{ $fieldName }}) FieldMask() *{{ nestedMessageName $field }}_FieldMask {
	return &{{ nestedMessageName $field }}_FieldMask{
		paths: nm.parent.fm.{{ $fmField.Name.UpperCamelCase }}.Paths,
		prefix: nm.path,
	}
}

		{{ end }}
	{{ end }}
{{ end }}

{{ end }}