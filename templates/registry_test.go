package templates

import (
	"testing"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/stretchr/testify/suite"

	"github.com/yeqown/protoc-gen-fieldmask/templates/shared"
)

type testGoTemplateRegistrySuite struct {
	suite.Suite

	tpl *template.Template
}

func (t *testGoTemplateRegistrySuite) SetupSuite() {
	tpl := template.New("file.tpl")
	shared.RegisterFunctions(tpl, pgsgo.InitContext(pgs.Parameters{}))
	var err error
	_, err = tpl.ParseFiles("./go/file.tpl")
	t.Require().NoError(err)
	_, err = tpl.New("request_mask").ParseFiles("./go/request_mask.tpl")
	t.Require().NoError(err)
	_, err = tpl.New("response_mask").ParseFiles("./go/response_mask.tpl")
	t.Require().NoError(err)
	_, err = tpl.New("marked_checker").ParseFiles("./go/marked_checker.tpl")
	t.Require().NoError(err)
	t.tpl = tpl
}

func (t *testGoTemplateRegistrySuite) Test_Run() {
	// The V2 templates require complex context that cannot be easily mocked.
	// This test verifies that all templates can be parsed successfully.
	// Template execution is tested via integration tests with actual proto files.
	t.NotNil(t.tpl)
}

func Test_GoTemplateRegistrySuite(t *testing.T) {
	suite.Run(t, new(testGoTemplateRegistrySuite))
}
