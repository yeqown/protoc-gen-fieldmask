package fieldmask_test

import (
	"testing"

	"github.com/yeqown/protoc-gen-fieldmask/proto/fieldmask"
	testdata "github.com/yeqown/protoc-gen-fieldmask/proto/fieldmask/testdata"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test_FieldMask_Masked(t *testing.T) {
	nfm := fieldmask.NewWithPaths(
		"a",         // basic type list field
		"a.invalid", // invalid field path
		"b.m1",      // message list field
		"c",         // basic field
		"e.e1",      // embedded message field
	)

	assert.True(t, nfm.Masked("a"))
	assert.True(t, nfm.Masked("a.invalid"))
	assert.True(t, nfm.Masked("b"))
	assert.True(t, nfm.Masked("b.m1"))
	assert.False(t, nfm.Masked("b.m2"))
	assert.True(t, nfm.Masked("c"))
	assert.False(t, nfm.Masked("d"))
	assert.True(t, nfm.Masked("e"))
	assert.True(t, nfm.Masked("e.e1"))
	assert.False(t, nfm.Masked("e.e2"))
}

func Test_FieldMask_Filter(t *testing.T) {
	nfm := fieldmask.NewWithPaths(
		"a",         // basic type list field
		"a.invalid", // invalid field path
		"b.m1",      // message list field
		"c",         // basic field
		"e.e1",      // embedded message field
	)

	// executing filter
	in := &testdata.T1{
		A: []int32{1, 2, 3},
		B: []*testdata.RepeatMessage{
			{
				M1: "b.m1",
				M2: 123,
			},
		},
		C: "t1.c",
		D: true,
		E: &testdata.T1_Embed{
			E1: 123,
			E2: "embed.e2",
		},
	}
	nfm.Filter(in)

	assert.NotEmpty(t, in.A)
	assert.NotEmpty(t, in.B)
	assert.NotEmpty(t, in.B[0].M1)
	assert.Empty(t, in.B[0].M2)
	assert.NotEmpty(t, in.C)
	assert.Empty(t, in.D)
	assert.NotEmpty(t, in.E)
	assert.NotEmpty(t, in.E.E1)
	assert.Empty(t, in.E.E2)

	byts, _ := protojson.Marshal(in)
	t.Logf("output: %s", byts)
}

func Test_FieldMask_Prune(t *testing.T) {
	nfm := fieldmask.NewWithPaths(
		"a",         // basic type list field
		"a.invalid", // invalid field path would mislead the prune action.
		"b.m1",      // message list field
		"c",         // basic field
		"e.e1",      // embedded message field
		"e.e2",      // embedded message field
	)

	// executing filter
	in := &testdata.T1{
		A: []int32{1, 2, 3},
		B: []*testdata.RepeatMessage{
			{
				M1: "b.m1",
				M2: 123,
			},
		},
		C: "t1.c",
		D: true,
		E: &testdata.T1_Embed{
			E1: 123,
			E2: "embed.e2",
		},
	}
	nfm.Prune(in)

	assert.Empty(t, in.A)
	assert.NotEmpty(t, in.B)
	assert.Empty(t, in.B[0].M1)
	assert.NotEmpty(t, in.B[0].M2)
	assert.Empty(t, in.C)
	assert.NotEmpty(t, in.D)
	assert.NotEmpty(t, in.E)
	assert.Empty(t, in.E.E1)
	assert.Empty(t, in.E.E2)

	byts, _ := protojson.Marshal(in)
	t.Logf("output: %s", byts)
}
