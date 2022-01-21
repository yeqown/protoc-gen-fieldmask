package templates

import (
	"os"
	"testing"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/stretchr/testify/suite"

	"github.com/yeqown/protoc-gen-fieldmask/internal/templates/shared"
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
	_, err = tpl.New("message").ParseFiles("./go/message.tpl")
	t.Require().NoError(err)
	t.tpl = tpl
}

func (t testGoTemplateRegistrySuite) Test_Run() {
	err := t.tpl.Execute(os.Stderr, map[string]interface{}{
		"Name": "test",
	})
	t.NoError(err)
}

func Test_GoTemplateRegistrySuite(t *testing.T) {
	suite.Run(t, new(testGoTemplateRegistrySuite))
}
