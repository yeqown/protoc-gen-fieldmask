package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/yeqown/protoc-gen-fieldmask/internal/module"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(pgs.DebugEnv("DEBUG_PGFM"), pgs.SupportedFeatures(&optional)).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
