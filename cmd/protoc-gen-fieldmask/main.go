package main

import (
	"github.com/yeqown/protoc-gen-fieldmask/pkg/module"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(pgs.DebugEnv("DEBUG_PGFM"), pgs.SupportedFeatures(&optional)).
		RegisterModule(module.FieldMask()).
		RegisterPostProcessor(
			pgsgo.GoFmt(),
		).
		Render()
}
