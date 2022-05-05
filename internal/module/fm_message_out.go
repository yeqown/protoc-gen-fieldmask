package module

import pgs "github.com/lyft/protoc-gen-star"

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
	*checkInMessageVO

	InMessage        pgs.Message
	OutMessage       pgs.Message
	GenOutMessageVar bool

	// OutMessagePkgName is the package name of the out message,
	// it's not empty while OutMessage is imported from other protobuf file.
	OutMessagePkgName string
}
