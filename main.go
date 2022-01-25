package main

import (
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/yeqown/protoc-gen-fieldmask/internal/module"

	pgs "github.com/lyft/protoc-gen-star"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(pgs.DebugEnv("DEBUG_PGFM"), pgs.SupportedFeatures(&optional)).
		RegisterModule(module.FieldMask(tplRegistryFactory)).
		RegisterPostProcessor(
			pgsgo.GoFmt(),
		).
		Render()
}
