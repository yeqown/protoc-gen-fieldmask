{{ $inMessageName := .InMessage.Name }}
{{ $outMessageName := .OutMessage.Name }}
{{ $fmField := .FieldMaskField }}

func (x *{{ $inMessageName }}) {{ $fmField.Name.UpperCamelCase }}_WithMode(mode pbfieldmask.MaskMode) *{{ $inMessageName }}_FieldMask {
    fm := &{{ $inMessageName }}_FieldMask{
        maskMode: mode,
        mask: pbfieldmask.New(x.{{$fmField.Name.UpperCamelCase}}),
    }

    return fm
}

// {{ $fmField.Name.UpperCamelCase }}_Filter generates *{{ $inMessageName }}_FieldMask with filter mode, so that
// only the fields in the {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }} will be
// appended into the {{ $inMessageName }}.
func (x *{{ $inMessageName }}) {{ $fmField.Name.UpperCamelCase }}_Filter() *{{ $inMessageName }}_FieldMask {
	return x.{{ $fmField.Name.UpperCamelCase }}_WithMode(pbfieldmask.MaskMode_Filter)
}

// {{ $fmField.Name.UpperCamelCase }}_Prune generates *{{ $inMessageName }}_FieldMask with prune mode, so that
// only the fields NOT in the {{ $inMessageName }}.{{ $fmField.Name.UpperCamelCase }} will be
// appended into the {{ $inMessageName }}.
func (x *{{ $inMessageName }}) {{ $fmField.Name.UpperCamelCase }}_Prune() *{{ $inMessageName }}_FieldMask {
	return x.{{ $fmField.Name.UpperCamelCase }}_WithMode(pbfieldmask.MaskMode_Prune)
}

// {{ $inMessageName }}_FieldMask provide provide helper functions to deal with FieldMask.
type {{ $inMessageName }}_FieldMask struct {
    maskMode pbfieldmask.MaskMode
    mask pbfieldmask.NestedFieldMask
}