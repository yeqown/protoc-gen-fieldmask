package module

import pgs "github.com/lyft/protoc-gen-star"

// FieldMask is a helper type for generating field masks.
type FieldMask struct{}

// New configures the module with an instance of FieldMask
func New() pgs.Module {
	return &FieldMask{}
}

func (f *FieldMask) Name() string {
	return "fieldmask"
}

func (f *FieldMask) InitContext(c pgs.BuildContext) {
	//TODO implement me
	panic("implement me")
}

func (f *FieldMask) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	//TODO implement me
	panic("implement me")
}
