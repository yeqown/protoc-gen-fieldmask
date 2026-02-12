package module

import (
	pgs "github.com/lyft/protoc-gen-star"

	fieldmask "github.com/yeqown/protoc-gen-fieldmask/proto/fieldmask"
)

type outFieldMaskContext struct {
	File           pgs.File
	FieldMaskPairs []fmMessagePair
	// ImportPaths is a slice of the OutMessage imported path and package name.
	ImportPaths []importPathPair
}

type importPathPair struct {
	ImportPath string
	// PkgName is the package name of the OutMessage,
	// it keeps the same as the fmMessagePair.OutMessagePkgName.
	PkgName string
}

type fmMessagePair struct {
	MethodOptions     *fieldmask.MethodOptions
	FieldMaskField    pgs.Field
	InMessage         pgs.Message
	OutMessage        pgs.Message
	OutMessagePkgName string
}

// checkMethodOptions checks if the method has field mask options configured
func checkMethodOptions(method pgs.Method, debugf func(string, ...interface{})) (*fieldmask.MethodOptions, bool) {
	if method == nil {
		return nil, false
	}

	var opts fieldmask.MethodOptions
	_, err := method.Extension(fieldmask.E_Rpc, &opts)
	if err != nil {
		return nil, false
	}

	debugf("method (%s) has fieldmask options", method.Name())
	return &opts, true
}

// findFieldMaskField finds the field mask field in the request message
func findFieldMaskField(message pgs.Message, fieldName string) (pgs.Field, bool) {
	if message == nil || fieldName == "" {
		return nil, false
	}

	fields := message.Fields()
	for _, field := range fields {
		if field.Name().String() == fieldName {
			// Check if it's a google.protobuf.FieldMask type
			if field.Type().ProtoType() == pgs.MessageT &&
				field.Descriptor().GetTypeName() == ".google.protobuf.FieldMask" {
				return field, true
			}
		}
	}

	return nil, false
}

// checkFieldOptions checks if a field has field mask options configured
func checkFieldOptions(field pgs.Field) (*fieldmask.FieldOptions, bool) {
	if field == nil {
		return nil, false
	}

	var opts fieldmask.FieldOptions
	_, err := field.Extension(fieldmask.E_Field, &opts)
	if err != nil {
		return nil, false
	}

	return &opts, true
}
