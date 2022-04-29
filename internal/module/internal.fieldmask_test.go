package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractPackagePrefix(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name            string
		args            args
		wantPkgPrefix   string
		wantMessageName string
	}{
		{
			name:            "normal",
			args:            args{name: "foo.bar.baz"},
			wantPkgPrefix:   "foo.bar",
			wantMessageName: "baz",
		},
		{
			name:            "normal",
			args:            args{name: "foo.bar.baz.qux"},
			wantPkgPrefix:   "foo.bar.baz",
			wantMessageName: "qux",
		},
		{
			name:            "normal",
			args:            args{name: "foo.bar"},
			wantPkgPrefix:   "foo",
			wantMessageName: "bar",
		},
		{
			name:            "normal",
			args:            args{name: "foo"},
			wantPkgPrefix:   "",
			wantMessageName: "foo",
		},
		{
			name:            "normal",
			args:            args{name: ""},
			wantPkgPrefix:   "",
			wantMessageName: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPkgPrefix, gotMessageName := extractPackagePrefix(tt.args.name)
			assert.Equal(t, tt.wantPkgPrefix, gotPkgPrefix)
			assert.Equal(t, tt.wantMessageName, gotMessageName)
		})
	}
}
