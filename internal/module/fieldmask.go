package module

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/yeqown/protoc-gen-fieldmask/internal/templates"
)

const (
	moduleName  = "fieldmask"
	langParam   = "lang"
	moduleParam = "module"
)

var (
	_ pgs.Module = (*FieldMaskModule)(nil)
)

// FieldMaskModule is a helper type for generating field masks.
type FieldMaskModule struct {
	*pgs.ModuleBase

	ctx pgsgo.Context

	registryFactory func(ctx pgsgo.Context) *templates.Registry
	registry        *templates.Registry
}

// FieldMask configures the module with an instance of FieldMaskModule
func FieldMask(registryFactory func(ctx pgsgo.Context) *templates.Registry) pgs.Module {
	return &FieldMaskModule{
		ModuleBase:      &pgs.ModuleBase{},
		registryFactory: registryFactory,
	}
}

func (m *FieldMaskModule) Name() string {
	return moduleName
}

func (m *FieldMaskModule) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
	m.registry = m.registryFactory(m.ctx)
}

func (m *FieldMaskModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	lang := m.Parameters().Str(langParam)
	m.Assert(lang != "", " `lang` parameter must be set")
	module := m.Parameters().Str(moduleParam)
	_ = module

	tpls := m.registry.Load(lang)
	m.Debugf("Loaded %d templates for %s", len(tpls), lang)
	// m.Assert(len(tpls) != 0, " could not find templates for `lang`: ", lang)
	_ = tpls

	outMessageVars := map[string]struct{}{}
	for _, f := range targets {
		m.Push(f.Name().String()).Debug("fieldmask")

		shouldGen := false // shouldGen indicates whether the file should be generated with '*.fm.go'.
		fileMessagesMapping := make(map[string]pgs.Message)
		paris := make([]fmMessagePair, 0, 2)
		for _, message := range f.AllMessages() {
			fileMessagesMapping[message.Name().String()] = message
			// _, _ = i, message
			// DONE(@yeqown): check message contains a field google.protobuf.FieldMask and
			// specify the fieldmask.option.Option as field option.
			r := m.checkInMessage(message)
			if r.invalid() {
				continue
			}

			m.Debugf("message %s has fieldmask field %s", message.Name(), r.FieldMaskField.Name())
			paris = append(paris, fmMessagePair{
				checkInMessageVO: r,
				InMessage:        message,
				OutMessage:       nil,
			})

			shouldGen = true
		}

		outCtx := &outFieldMaskContext{
			File:           f,
			FieldMaskPairs: make([]fmMessagePair, 0, len(paris)),
		}

		if shouldGen {
			m.Debugf("file (%s) is planned to generate user.pb.fm.go", f.Name().String())
			// 从当前文件中找到匹配的 message，如果没有找到，则跳过
			for idx, pair := range paris {
				ok := false
				if pair.OutMessage == nil {
					outMessageName := pair.checkInMessageVO.FieldMaskOption.GetOut().GetMessage()
					paris[idx].OutMessage, ok = fileMessagesMapping[outMessageName]
					if !ok {
						continue
					}
					if _, found := outMessageVars[outMessageName]; !found {
						paris[idx].GenOutMessageVar = true
						outMessageVars[outMessageName] = struct{}{}
					}
				}

				outCtx.FieldMaskPairs = append(outCtx.FieldMaskPairs, paris[idx])
			}

			filename := m.ctx.OutputPath(f).SetExt(".fm.go")
			for idx := range tpls {
				m.Debugf("add template %s to %s", tpls[idx].Name(), filename.String())
				m.AddGeneratorTemplateFile(filename.String(), tpls[idx], outCtx)
			}
		}

		m.Pop()
	}

	return m.Artifacts()
}
