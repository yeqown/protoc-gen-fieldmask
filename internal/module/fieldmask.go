package module

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"strings"

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

	// pkgMessageCache map[fullPathMessage]Message eg. google.protobuf.Timestamp: Timestamp
	pkgMessageCache *pkgMessageCache
}

// FieldMask configures the module with an instance of FieldMaskModule
func FieldMask(registryFactory func(ctx pgsgo.Context) *templates.Registry) pgs.Module {
	return &FieldMaskModule{
		ModuleBase:      &pgs.ModuleBase{},
		ctx:             nil,
		registryFactory: registryFactory,
		registry:        nil,
		lang:            "",
		pkgMessageCache: newCache(0),
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
	for filename, f := range targets {
		_ = filename
		m.Push(f.Name().String()).Debug("fieldmask")

		pairs, mm := m.parse(f)
		if len(pairs) <= 0 {
			m.Pop()
			continue
		}

		data := &outFieldMaskContext{
			File:           f,
			FieldMaskPairs: pairs,
		}

		m.generate(data, mm, packages)
		m.Pop()
	}

	return m.Artifacts()
}

// parse the given file and returns a list of messages that
// contain a fieldmask field.
func (m *FieldMaskModule) parse(target pgs.File) (pairs []fmMessagePair, mm map[string]pgs.Message) {
	mm = make(map[string]pgs.Message, 8)
	paris := make([]fmMessagePair, 0, 2)
	for _, message := range target.AllMessages() {
		mm[message.Name().String()] = message

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

	return paris, mm
}

// locateMessage finds the message that specified by InMessage. Firstly, it
// judges whether the message is current file's message. If not, it will try to
// find the message in import packages.
func (m *FieldMaskModule) locateMessage(
	name string, mm map[string]pgs.Message, packages map[string]pgs.Package) (message pgs.Message, ok bool) {
	if name == "" {
		m.Debug("locateMessage: message name is empty")
		return nil, false
	}

	if pkg, messageName := extractPackagePrefix(name); pkg != "" {
		// if the pkg is qualified, it means the message is in import packages.
		message, ok = lookupMessageFromPackageCached(packages, m.pkgMessageCache, pkg, messageName)
		if ok {
			return message, true
		}
	}

	// if the pkg is not qualified, it means the message is in current file.
	message, ok = mm[name]
	return message, ok
}

// lookupMessageFromPackageCached finds the message from the given package and cache.
func lookupMessageFromPackageCached(
	packages map[string]pgs.Package, cache *pkgMessageCache, pkgName, messageName string) (message pgs.Message, ok bool) {
	if pkgName == "" {
		return nil, false
	}

	fqn := "." + pkgName + "." + messageName
	if message, ok = cache.cached(fqn); ok {
		return message, ok
	}

	// locate target package, if not found return nil.
	pkg, hit := packages[pkgName]
	if !hit {
		return nil, false
	}

	// find the message from the package.
	files := pkg.Files()
	for i := 0; i < len(files); i++ {
		found := false
		_file := files[i]
		// if _file has been parsed, just skip current _file.
		if cache.isFileParsed(_file.InputPath().String()) {
			continue
		}

		// caching message in current _file, if the message has been found,
		// do not return until all messages in current _file are cached.
		for idx, _m := range _file.Messages() {
			cache.cache(_m.FullyQualifiedName(), _file.Messages()[idx])
			if _m.Name().String() == messageName {
				found = true
				ok = true
				message = _m
			}
		}

		cache.markFileParsed(_file.InputPath().String())
		if found {
			break
		}
	}

	return message, ok
}

// extractPackagePrefix extracts the package prefix from the given message type name.
// e.g.
// extractPackagePrefix("com.pkg.Message") => "com.pkg" "Message"
func extractPackagePrefix(name string) (pkgPrefix, messageName string) {
	q := strings.Split(name, ".")
	switch c := len(q); c {
	case 0, 1:
		return "", name
	default:
		return strings.Join(q[:c-1], "."), q[c-1]
	}
}

// generate works in file domain, and generate the fieldmask files with templates.
// It will generate the fieldmask files with the given data.
func (m *FieldMaskModule) generate(data *outFieldMaskContext, mm map[string]pgs.Message, packages map[string]pgs.Package) {
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
		if msg, found := m.locateMessage(outMessageName, mm, packages); !found {
			m.Debugf("message %s is not found", outMessageName)
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
