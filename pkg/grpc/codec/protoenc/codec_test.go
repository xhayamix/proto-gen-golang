package protoenc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/QualiArts/campus-server/pkg/domain/proto/client/api"
	"github.com/QualiArts/campus-server/pkg/grpc/codec/protoenc/testdata"
)

//nolint:gochecknoinits // grpc/encodingがpackage変数をスレッドセーフではなく更新するためtestではinitで行う
func init() {
	Register([]byte("test"))
}

func TestCodec_Marshal(t *testing.T) {
	c := Codec{}

	data, err := c.Marshal(&pb.HealthCheckRequest{Service: "test"})
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)
	messageData := data[headerLengthBytes:]
	compress := int(messageData[0])
	assert.Equal(t, 0, compress)

	var v pb.HealthCheckRequest
	err = c.Unmarshal(data, &v)
	assert.NoError(t, err)
	assert.Equal(t, "test", v.Service)
}

func TestCodec_MarshalBigValue(t *testing.T) {
	c := Codec{}

	bigText := string(make([]byte, compressThresholdBytes))
	data, err := c.Marshal(&pb.HealthCheckRequest{Service: bigText})
	if assert.NoError(t, err) {
		assert.True(t, len(data) > 0)
		messageData := data[headerLengthBytes:]
		compress := int(messageData[0])
		assert.Equal(t, 1, compress)
	}

	var v pb.HealthCheckRequest
	err = c.Unmarshal(data, &v)
	assert.NoError(t, err)
	assert.Equal(t, bigText, v.Service)
}

func TestCodec_MarshalZeroValue(t *testing.T) {
	c := Codec{}

	// 空のデータを使った Unmarshal -> Marshal テスト
	var v pb.HealthCheckRequest
	err := c.Unmarshal([]byte{}, &v)
	assert.NoError(t, err)
	assert.Equal(t, "", v.Service)

	data, err := c.Marshal(&v)
	assert.NoError(t, err)
	assert.True(t, len(data) == 0)
}

func TestCodec_MarshalEmpty(t *testing.T) {
	c := Codec{}

	// Marshal -> Unmarshal テスト
	data, err := c.Marshal(&emptypb.Empty{})
	assert.NoError(t, err)
	assert.Empty(t, data)

	var v emptypb.Empty
	err = c.Unmarshal(data, &v)
	assert.NoError(t, err)
}

func TestCodec_UnmarshalRawFields(t *testing.T) {
	// version間で差分のある定義の中でも元のbyte配列に戻すことが可能かどうかテスト
	// 暗号化のテストではないのでcodec関係なくprotobufの検証

	v1 := &testdata.Version1{
		Name1: "v1",
		Name2: "unknownField",
	}
	v1Buf, err := marshal(v1)
	assert.NoError(t, err)

	v2 := new(testdata.Version2)
	err = unmarshal(v1Buf, v2)
	assert.NoError(t, err)

	v2Buf, err := marshal(v2)
	assert.NoError(t, err)

	assert.Equal(t, v1Buf, v2Buf)
}
