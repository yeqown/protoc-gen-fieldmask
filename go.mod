module github.com/yeqown/protoc-gen-fieldmask

go 1.25

require (
	github.com/iancoleman/strcase v0.2.0
	github.com/lyft/protoc-gen-star v0.6.2
	github.com/stretchr/testify v1.11.1
	github.com/yeqown/protoc-gen-fieldmask/protobuf v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.11
)

replace github.com/yeqown/protoc-gen-fieldmask/protobuf => ./protobuf

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
