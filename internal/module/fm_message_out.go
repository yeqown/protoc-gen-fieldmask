package module

import pgs "github.com/lyft/protoc-gen-star"

type outFieldMaskContext struct {
	File           pgs.File
	FieldMaskPairs []fmMessagePair
}

type fmMessagePair struct {
	*checkInMessageVO

	FieldMaskField pgs.Field
	InMessage      pgs.Message
	OutMessage     pgs.Message
}
