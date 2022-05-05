package module

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"regexp"
	"strconv"
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

		ctx := &outFieldMaskContext{
			File:           f,
			FieldMaskPairs: pairs,
			ImportPaths:    make([]importPathPair, 0, 4),
		}

		// collect import paths and consummate the out messages.
		m.consummate(ctx, mm, packages)

		m.generate(ctx)
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
// find the message in import packages, at the same time, the target package import path
// and package name will be returned.
func (m *FieldMaskModule) locateMessage(
	name string, mm map[string]pgs.Message, packages map[string]pgs.Package,
) (message pgs.Message, importPath, packageName string, ok bool) {
	if name == "" {
		m.Debug("locateMessage: message name is empty")
		return nil, "", "", false
	}

	if pkg, messageName := extractPackagePrefix(name); pkg != "" {
		// if the pkg is qualified, it means the message is in import packages.
		message, ok = lookupMessageFromPackageCached(packages, m.pkgMessageCache, pkg, messageName)
		if ok {
			switch m.lang {
			case "go":
				option := message.File().Descriptor().GetOptions().GetGoPackage()
				importPath, packageName = resolveGoPackageOption(option)
				// TODO(@yeqown): support multi language. now only support go.
			}

			return message, importPath, packageName, true
		}
	}

	// if the pkg is not qualified, it means the message is in current file.
	message, ok = mm[name]
	return message, "", "", ok
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

var nonAlphaNumPattern = regexp.MustCompile("[^a-zA-Z0-9]")

// parseGoPackageOption parses the go_package option from the option string.
//
// .eg1.
// go_package="example.com/foo/bar;baz" should have a package name of `baz`
// and an import path of `example.com/foo/bar`.
// .eg2.
// go_package="example.com/foo/bar" should have a package name of `bar`
// and an import path of `example.com/foo/bar`.
func resolveGoPackageOption(option string) (path, pkg string) {
	if option == "" {
		return "", ""
	}

	// .eg1: example.com/foo/bar;baz
	if idx := strings.LastIndex(option, ";"); idx > -1 {
		path = option[:idx]
		pkg = nonAlphaNumPattern.ReplaceAllString(option[idx+1:], "_")
		return
	}

	// .eg2: example.com/foo/bar
	if idx := strings.LastIndex(option, "/"); idx > -1 {
		path = option
		pkg = nonAlphaNumPattern.ReplaceAllString(option[idx+1:], "_")
		return
	}

	return "", ""
}

// extractPackagePrefix extracts the package prefix from the given message type name.
// e.g.
// extractPackagePrefix("com.pkg.Message") => "com.pkg", "Message"
// extractPackagePrefix("Message") => "", "Message" which means the message is in the current file.
func extractPackagePrefix(name string) (pkgPrefix, messageName string) {
	q := strings.Split(name, ".")
	switch c := len(q); c {
	case 0, 1:
		return "", name
	default:
		return strings.Join(q[:c-1], "."), q[c-1]
	}
}

// consummate fm pairs with full qualified OutMessage which means it has
// import path and package name as long as it is one message type defined
// in another protobuf file.
func (m *FieldMaskModule) consummate(
	ctx *outFieldMaskContext, mm map[string]pgs.Message, packages map[string]pgs.Package) {
	m.Debugf("consummating fm pairs with full qualified OutMessage")

	// uniq the out message initialization statement
	uniq := make(map[string]struct{}, len(ctx.FieldMaskPairs))
	// uniqImportPath the out message import statement, map[importPath]packageName.
	uniqImportPath := make(map[string]string, len(ctx.FieldMaskPairs))
	// uniqPkgAlias the out message package alias, map[packageName]count
	uniqPkgAlias := make(map[string]uint8, len(ctx.FieldMaskPairs))
	filteredFmPairs := make([]fmMessagePair, 0, len(ctx.FieldMaskPairs))

	// consummate outFieldMaskContext.FieldMaskPairs and outFieldMaskContext.ImportPaths.
	for idx, pair := range ctx.FieldMaskPairs {
		if pair.OutMessage != nil {
			filteredFmPairs = append(filteredFmPairs, ctx.FieldMaskPairs[idx])
			continue
		}

		outMessageName := pair.checkInMessageVO.FieldMaskOption.GetOut().GetMessage()
		msg, importPath, pkgName, found := m.locateMessage(outMessageName, mm, packages)
		if !found {
			m.Debugf("message %s is not found", outMessageName)
			continue
		}
		ctx.FieldMaskPairs[idx].OutMessage = msg

		// only the OutMessage is imported from third protobuf file.
		if importPath != "" && pkgName != "" {
			// DONE(@yeqown): let import paths unique in same file, package names unique for the same package name.
			if c, ok := uniqPkgAlias[pkgName]; ok {
				if c != 0 {
					pkgName = pkgName + strconv.Itoa(int(uniqPkgAlias[pkgName]))
				}
				uniqPkgAlias[pkgName]++
			}
			if _, ok := uniqImportPath[importPath]; !ok {
				uniqImportPath[importPath] = pkgName
				ctx.ImportPaths = append(ctx.ImportPaths, importPathPair{
					ImportPath: importPath,
					PkgName:    pkgName,
				})
			}
			// DONE(@yeqown): use importPath and pkgName to generate template.
			ctx.FieldMaskPairs[idx].OutMessagePkgName = pkgName
		}

		// FIXED(@vaidasn): Generate out message vars only once per type #8
		if _, dup := uniq[outMessageName]; !dup {
			uniq[outMessageName] = struct{}{}
			ctx.FieldMaskPairs[idx].GenOutMessageVar = true
		}
		filteredFmPairs = append(filteredFmPairs, ctx.FieldMaskPairs[idx])
	}
	ctx.FieldMaskPairs = filteredFmPairs
}

// generate works in file domain, and generate the fieldmask files with templates.
// It will generate the fieldmask files with the given data.
func (m *FieldMaskModule) generate(data *outFieldMaskContext) {
	m.Debugf("file (%s) is planned to generate user.pb.fm.go", data.File.Name().String())

	setting := m.registry.Load(m.lang)
	m.Debugf("Loaded %d templates for %s", len(setting.Templates), m.lang)
	filename := m.ctx.
		OutputPath(data.File).
		SetExt(setting.Ext)
	for _, tpl := range setting.Templates {
		m.Debugf("add template %s to %s", tpl.Name(), filename.String())
		m.AddGeneratorTemplateFile(filename.String(), tpl, data)
	}
}
