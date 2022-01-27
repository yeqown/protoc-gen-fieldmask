package module_test

import (
	"bytes"
	"os"
	"testing"

	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/yeqown/protoc-gen-fieldmask/internal/module"
	"github.com/yeqown/protoc-gen-fieldmask/internal/templates"

	pgs "github.com/lyft/protoc-gen-star"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type moduleTestSuite struct {
	suite.Suite

	module pgs.Module
}

func (m moduleTestSuite) TearDownSuite() {
	// pass
}

// SetupSuite runs once, before all tests in the suite
// https://pkg.go.dev/github.com/lyft/protoc-gen-star/protoc-gen-debug#section-readme
func (m *moduleTestSuite) SetupSuite() {
	registryFactory := func(ctx pgsgo.Context) *templates.Registry {
		registry := templates.New()
		registry.Register("go", nil)
		return registry
	}

	m.module = module.FieldMask(registryFactory)
}

func mutateLangParam(lang string) pgs.ParamMutator {
	return func(p pgs.Parameters) {
		p["lang"] = lang
	}
}

func (m *moduleTestSuite) Test_ForDebug() {
	// please look up the README at repository root directory to see how to
	// generate the `testdata` and code_generator_request binary.
	req, err := os.Open("./testdata/code_generator_request.pb.bin")
	m.NoError(err)

	fs := afero.NewMemMapFs()
	res := &bytes.Buffer{}

	pgs.Init(
		pgs.ProtocInput(req),                    // use the pre-generated request
		pgs.ProtocOutput(res),                   // capture CodeGeneratorResponse
		pgs.FileSystem(fs),                      // capture any custom files written directly to disk
		pgs.MutateParams(mutateLangParam("go")), // mutate params
	).
		RegisterModule(m.module).
		Render()
}

func Test_module(t *testing.T) {
	suite.Run(t, new(moduleTestSuite))
}
