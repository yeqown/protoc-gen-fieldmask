package module

import pgs "github.com/lyft/protoc-gen-star"

// FieldMasker is a helper type for generating field masks.
type FieldMasker struct{}

func New() *FieldMasker {
	return &FieldMasker{}
}

func (f FieldMasker) Name() string {
	//TODO implement me
	panic("implement me")
}

func (f FieldMasker) InitContext(c pgs.BuildContext) {
	//TODO implement me
	panic("implement me")
}

func (f FieldMasker) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	//TODO implement me
	panic("implement me")
}
