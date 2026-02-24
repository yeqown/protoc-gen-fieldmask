package module

import (
	pgs "github.com/lyft/protoc-gen-star"

	"github.com/yeqown/protoc-gen-fieldmask/protobuf"
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
	Method            pgs.Method
	MethodOptions     *protobuf.MethodOptions
	FieldMaskField    pgs.Field
	InMessage         pgs.Message
	OutMessage        pgs.Message
	OutMessagePkgName string
}

// checkMethodOptions checks if the method has field mask options configured
func checkMethodOptions(method pgs.Method, debugf func(string, ...interface{})) (*protobuf.MethodOptions, bool) {
	if method == nil {
		return nil, false
	}

	var opts protobuf.MethodOptions
	_, err := method.Extension(protobuf.E_Rpc, &opts)
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
