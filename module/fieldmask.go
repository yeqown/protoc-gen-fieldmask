package module

import (
	"regexp"
	"strconv"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	templates2 "github.com/yeqown/protoc-gen-fieldmask/templates"
)

var goPkgNamePattern = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

var invalidPkgCharsPattern = regexp.MustCompile("[^a-zA-Z0-9]")

const (
	moduleName  = "fieldmask"
	langParam   = "lang"
	moduleParam = "module"
)

var (
	_ pgs.Module = (*protocGenFieldmask)(nil)
)

// protocGenFieldmask is a helper type for generating field masks.
type protocGenFieldmask struct {
	*pgs.ModuleBase

	ctx pgsgo.Context

	registryFactory func(ctx pgsgo.Context) *templates2.Registry
	registry        *templates2.Registry
	lang            string

	// pkgMessageCache map[fullPathMessage]Message eg. google.protobuf.Timestamp: Timestamp
	pkgMessageCache *pkgMessageCache
}

// NewModule configures the module with an instance of protocGenFieldmask
func NewModule() pgs.Module {
	return &protocGenFieldmask{
		ModuleBase:      &pgs.ModuleBase{},
		ctx:             nil,
		registryFactory: templates2.RegistryFactory,
		registry:        nil,
		lang:            "",
		pkgMessageCache: newCache(0),
	}
}

func (m *protocGenFieldmask) Name() string {
	return moduleName
}

func (m *protocGenFieldmask) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
	m.registry = m.registryFactory(m.ctx)
}

func (m *protocGenFieldmask) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	m.lang = m.Parameters().Str(langParam)
	m.Assert(m.lang != "", " `lang` parameter must be set")
	module := m.Parameters().Str(moduleParam)
	_ = module

	// range all target files, to locate RPC methods that have field mask options
	for filename, f := range targets {
		_ = filename
		m.Push(f.Name().String()).Debug("fieldmask")

		pairs := m.parseServices(f)
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
		m.consummate(ctx, packages)

		m.generate(ctx)
		m.Pop()
	}

	return m.Artifacts()
}

// parseServices parses the given file and returns a list of RPC methods that have field mask options
func (m *protocGenFieldmask) parseServices(target pgs.File) (pairs []fmMessagePair) {
	pairs = make([]fmMessagePair, 0, 2)

	// First collect all messages in the file
	messages := make(map[string]pgs.Message)
	for _, message := range target.AllMessages() {
		messages[message.Name().String()] = message
	}

	// Parse services and their RPC methods
	for _, service := range target.Services() {
		for _, method := range service.Methods() {
			methodOpts, found := checkMethodOptions(method, m.Debugf)
			if !found {
				continue
			}

			requestMessage := method.Input()
			responseMessage := method.Output()

			// Find the field mask field in the request message
			fieldMaskField, found := findFieldMaskField(requestMessage, methodOpts.FieldName)
			if !found {
				m.Debugf("field mask field '%s' not found in request message %s", methodOpts.FieldName, requestMessage.Name())
				continue
			}

			m.Debugf("method %s.%s has fieldmask field %s", service.Name(), method.Name(), fieldMaskField.Name())
			pairs = append(pairs, fmMessagePair{
				Method:            method,
				MethodOptions:     methodOpts,
				FieldMaskField:    fieldMaskField,
				InMessage:         requestMessage,
				OutMessage:        responseMessage,
				OutMessagePkgName: "",
			})
		}
	}

	return pairs
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
	// Find the last separator: either "/" or ";"
	// The LAST separator determines where the package name starts
	sepIdx := strings.LastIndexAny(option, "/;")
	if sepIdx == -1 {
		return "", ""
	}

	// Extract package name (after last separator)
	pkg = option[sepIdx+1:]

	// If the last separator is ";", everything before it is the path
	// If the last separator is "/", the entire option is the path
	if option[sepIdx] == ';' {
		path = option[:sepIdx]
	} else {
		path = option
	}

	// Clean package name by replacing non-alphanumeric characters with underscore
	pkg = nonAlphaNumPattern.ReplaceAllString(pkg, "_")

	return path, pkg
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
func (m *protocGenFieldmask) consummate(
	ctx *outFieldMaskContext, packages map[string]pgs.Package) {
	m.Debugf("consummating fm pairs with full qualified OutMessage")

	// uniqImportPath the out message import statement, map[importPath]packageName.
	uniqImportPath := make(map[string]string, len(ctx.FieldMaskPairs))
	// uniqPkgAlias the out message package alias, map[packageName]count
	uniqPkgAlias := make(map[string]uint8, len(ctx.FieldMaskPairs))

	// consummate outFieldMaskContext.FieldMaskPairs and outFieldMaskContext.ImportPaths.
	for idx, pair := range ctx.FieldMaskPairs {
		// Check if response message needs import resolution
		if pair.OutMessage == nil {
			continue
		}

		// Check if the response message is from another package
		// Skip if OutMessage is from same file as InMessage (no import needed)
		if pair.OutMessage.File() == pair.InMessage.File() {
			continue
		}
		if importPath, pkgName := m.getImportInfo(pair.OutMessage, packages); importPath != "" && pkgName != "" {
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
	}
}

// getImportInfo gets the import path and package name for a message
func (m *protocGenFieldmask) getImportInfo(message pgs.Message, packages map[string]pgs.Package) (importPath, pkgName string) {
	// Check if the message is from another file
	if message.File() == nil {
		return "", ""
	}

	// For now, only handle Go imports
	switch m.lang {
	case "go":
		option := message.File().Descriptor().GetOptions().GetGoPackage()
		importPath, pkgName = resolveGoPackageOption(option)
	}

	return importPath, pkgName
}

// generate works in file domain, and generate the fieldmask files with templates.
// It will generate the fieldmask files with the given data.
func (m *protocGenFieldmask) generate(data *outFieldMaskContext) {
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
