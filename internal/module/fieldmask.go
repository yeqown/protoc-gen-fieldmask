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
	lang            string
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
	m.lang = m.Parameters().Str(langParam)
	m.Assert(m.lang != "", " `lang` parameter must be set")
	module := m.Parameters().Str(moduleParam)
	_ = module

	// range all target files, to locate the target message which is specified
	// to generate field mask, and locate it's related OutMessage in current file and
	// imports, finally generate the fieldmask files with templates.
	for _, f := range targets {
		m.Push(f.Name().String()).Debug("fieldmask")

		pairs, mmapping := m.parse(f)
		if len(pairs) <= 0 {
			m.Pop()
			continue
		}

		data := &outFieldMaskContext{
			File:           f,
			FieldMaskPairs: pairs,
		}

		m.generate(data, mmapping)
		m.Pop()
	}

	return m.Artifacts()
}

// parse the given file and returns a list of messages that
// contain a fieldmask field.
func (m *FieldMaskModule) parse(target pgs.File) (pairs []fmMessagePair, mmapping map[string]pgs.Message) {
	mmapping = make(map[string]pgs.Message, 8)
	paris := make([]fmMessagePair, 0, 2)
	for _, message := range target.AllMessages() {
		mmapping[message.Name().String()] = message

		// DONE(@yeqown): check message contains a field google.protobuf.FieldMask and
		// specify the fieldmask.option.Option as field option.
		r := checkInMessage(message, m.Debugf)
		if r.invalid() {
			continue
		}

		m.Debugf("message %s has fieldmask field %s", message.Name(), r.FieldMaskField.Name())
		paris = append(paris, fmMessagePair{
			checkInMessageVO: r,
			InMessage:        message,
			OutMessage:       nil,
		})
	}

	return paris, mmapping
}

// locateMessage finds the message that specified by InMessage. Firstly, it
// judges whether the message is current file's message. If not, it will try to
// find the message in import packages.
func (m *FieldMaskModule) locateMessage(name string, mmaping map[string]pgs.Message) (pgs.Message, bool) {
	if name == "" {
		m.Debug("locateMessage: message name is empty")
		return nil, false
	}

	message, ok := mmaping[name]
	return message, ok

	// TODO(@yeqown): support different file's message.
	//return nil, false
}

// generate works in file domain, and generate the fieldmask files with templates.
// It will generate the fieldmask files with the given data.
func (m *FieldMaskModule) generate(data *outFieldMaskContext, mmaping map[string]pgs.Message) {
	m.Debugf("file (%s) is planned to generate user.pb.fm.go", data.File.Name().String())
	outMessageVars := make(map[string]struct{}, len(data.FieldMaskPairs))
	data2 := &outFieldMaskContext{
		File:           data.File,
		FieldMaskPairs: make([]fmMessagePair, 0, len(data.FieldMaskPairs)),
	}
	for idx, pair := range data.FieldMaskPairs {
		if pair.OutMessage != nil {
			data2.FieldMaskPairs = append(data2.FieldMaskPairs, data.FieldMaskPairs[idx])
			continue
		}

		outMessageName := pair.checkInMessageVO.FieldMaskOption.GetOut().GetMessage()
		if msg, found := m.locateMessage(outMessageName, mmaping); !found {
			m.Debugf("message %s is not found", msg.Name())
			continue
		} else {
			data.FieldMaskPairs[idx].OutMessage = msg
		}

		// FIXED(@vaidasn): Generate out message vars only once per type #8
		if _, found := outMessageVars[outMessageName]; !found {
			outMessageVars[outMessageName] = struct{}{}
			data.FieldMaskPairs[idx].GenOutMessageVar = true
		}

		data2.FieldMaskPairs = append(data2.FieldMaskPairs, data.FieldMaskPairs[idx])
	}

	setting := m.registry.Load(m.lang)
	m.Debugf("Loaded %d templates for %s", len(setting.Templates), m.lang)
	filename := m.ctx.
		OutputPath(data.File).
		SetExt(setting.Ext)
	for _, tpl := range setting.Templates {
		m.Debugf("add template %s to %s", tpl.Name(), filename.String())
		m.AddGeneratorTemplateFile(filename.String(), tpl, data2)
	}
}
