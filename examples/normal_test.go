package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"

	normal "examples/pb"
)

func Test_FieldMask_Filter(t *testing.T) {
	req := &normal.UserInfoRequest{
		UserId:    "123123",
		FieldMask: nil,
	}

	req.MaskOut_Email()
	req.MaskOut_Name()
	assert.Equal(t, 2, len(req.FieldMask.GetPaths()))

	byts, err := protojson.Marshal(req)
	require.NoError(t, err)
	t.Logf("req: %s", byts)

	filter := req.FieldMask_Filter()
	assert.NotNil(t, filter)

	resp := &normal.UserInfoResponse{
		UserId: "69781",
		Name:   "yeqown",
		Email:  "yeqown@gmail.com",
		Address: &normal.Address{
			Country:  "China",
			Province: "Sichuan",
		},
	}

	filter.Mask(resp)

	assert.NotEmpty(t, resp.Name)
	assert.NotEmpty(t, resp.Email)
	assert.Empty(t, resp.Address)
	assert.Empty(t, resp.UserId)

	byts2, err2 := protojson.Marshal(resp)
	require.NoError(t, err2)
	t.Logf("resp: %s", byts2)
}

func Test_FieldMask_Deep(t *testing.T) {
	req := &normal.UserInfoRequest{
		UserId:    "123123",
		FieldMask: nil,
	}

	req.MaskOut_Address_Country()
	assert.Equal(t, 1, len(req.FieldMask.GetPaths()))

	byts, err := protojson.Marshal(req)
	require.NoError(t, err)
	t.Logf("req: %s", byts)

	filter := req.FieldMask_Filter()
	assert.NotNil(t, filter)

	resp := &normal.UserInfoResponse{
		UserId: "69781",
		Name:   "yeqown",
		Email:  "yeqown@gmail.com",
		Address: &normal.Address{
			Country:  "China",
			Province: "Sichuan",
		},
	}

	filter.Mask(resp)

	assert.Empty(t, resp.Name)
	assert.Empty(t, resp.Email)
	assert.Empty(t, resp.UserId)
	assert.NotEmpty(t, resp.Address)
	assert.NotEmpty(t, resp.Address.Country)
	assert.Empty(t, resp.Address.Province)

	byts2, err2 := protojson.Marshal(resp)
	require.NoError(t, err2)
	t.Logf("resp: %s", byts2)
}

func Test_FieldMask_Deep2(t *testing.T) {
	req := &normal.UserInfoRequest{
		UserId:    "123123",
		FieldMask: nil,
	}

	req.MaskOut_Address_Country()
	assert.Equal(t, 1, len(req.FieldMask.GetPaths()))

	prune := req.FieldMask_Prune()
	assert.NotNil(t, prune)

	resp := &normal.UserInfoResponse{
		UserId: "69781",
		Name:   "yeqown",
		Email:  "yeqown@gmail.com",
		Address: &normal.Address{
			Country:  "China",
			Province: "Sichuan",
		},
	}

	prune.Mask(resp)

	assert.NotEmpty(t, resp.Name)
	assert.NotEmpty(t, resp.Email)
	assert.NotEmpty(t, resp.UserId)
	assert.NotEmpty(t, resp.Address)
	assert.NotEmpty(t, resp.GetAddress().GetProvince())
	assert.Empty(t, resp.GetAddress().GetCountry())

	byts2, err2 := protojson.Marshal(resp)
	require.NoError(t, err2)
	t.Logf("resp: %s", byts2)
}

func Test_FieldMask_Prune(t *testing.T) {
	req := &normal.UserInfoRequest{
		UserId:    "123123",
		FieldMask: nil,
	}

	req.MaskOut_Email()
	req.MaskOut_Name()
	assert.Equal(t, 2, len(req.FieldMask.GetPaths()))

	byts, err := protojson.Marshal(req)
	require.NoError(t, err)
	t.Logf("req: %s", byts)

	prune := req.FieldMask_Prune()
	assert.NotNil(t, prune)

	resp := &normal.UserInfoResponse{
		UserId: "69781",
		Name:   "yeqown",
		Email:  "yeqown@gmail.com",
		Address: &normal.Address{
			Country:  "China",
			Province: "Sichuan",
		},
	}

	prune.Mask(resp)

	assert.Empty(t, resp.Name)
	assert.Empty(t, resp.Email)
	assert.NotEmpty(t, resp.Address)
	assert.NotEmpty(t, resp.UserId)

	byts2, err2 := protojson.Marshal(resp)
	require.NoError(t, err2)
	t.Logf("resp: %s", byts2)
}

func Test_FieldMask_Masked(t *testing.T) {
	req := &normal.UserInfoRequest{
		UserId:    "123123",
		FieldMask: nil,
	}

	req.MaskOut_Email()
	req.MaskOut_Address_Country()
	assert.Equal(t, 2, len(req.FieldMask.GetPaths()))

	filter := req.FieldMask_Filter()
	assert.True(t, filter.MaskedOut_Address())
	assert.True(t, filter.MaskedOut_Address_Country())
	assert.True(t, filter.MaskedOut_Email())

	assert.False(t, filter.MaskedIn_UserId())
	assert.False(t, filter.MaskedOut_Name())
	assert.False(t, filter.MaskedOut_UserId())
	assert.False(t, filter.MaskedOut_Address_Province())
}
