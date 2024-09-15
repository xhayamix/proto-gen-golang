// 通常のproto codecに暗号化・圧縮処理を差し込むため
// https://github.com/grpc/grpc-go/blob/master/encoding/proto/proto.go
// のソースコードをベースに処理を追加している.
// 圧縮にはgzipを使うため暗号化前に圧縮処理が必要。そのためcodec後に動くcompressorでは使えない。

package protoenc

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/binary"
	"math/big"

	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/logging/app"
	"github.com/xhayamix/proto-gen-golang/pkg/util/cipher"
	"github.com/xhayamix/proto-gen-golang/pkg/util/randstr"
)

const (
	// 暗号化key文字最大数
	maxKeyLength = 32
	// バイナリデータのヘッダサイズ情報バイト数
	headerLengthBytes = 2
	// バイナリデータのメタデータバイト数
	metadataBytes = 2
	// 圧縮閾値Bytes
	compressThresholdBytes = 500000
)

// 暗号化Secret
var securitySecret []byte

const Name = "proto-enc"

// must be thread safe
func Register(secret []byte) {
	securitySecret = secret
	encoding.RegisterCodec(Codec{})
}

// codec implements gRPC codec interface
// that implementations of this interface must be thread safe; a codec's
// methods can be called from concurrent goroutines.
type Codec struct{}

// proto marshal
func marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

// proto unmarshal
func unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	// proto marshal
	out, err := marshal(v)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	// バイナリ変換の結果、データサイズがゼロだった場合はカスタムヘッダを付与しない
	if len(out) == 0 {
		return out, nil
	}

	// Messageの圧縮(サイズがしきい値を超えているときのみ圧縮)
	var compress int
	if len(out) >= compressThresholdBytes {
		buffer := new(bytes.Buffer)
		writer := gzip.NewWriter(buffer)
		if _, err = writer.Write(out); err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		if err := writer.Close(); err != nil {
			// エラーログを吐いて処理は進める
			app.GetLogger().Error(context.Background(), err.Error())
		}
		out = buffer.Bytes()
		compress = 1
	}

	// Message部分の暗号化
	maxInt := new(big.Int)
	maxInt.SetInt64(int64(maxKeyLength - 1))
	randomLength, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	keyBase, err := randstr.RandomString(1 + int(randomLength.Int64()))
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	cryptKey := []byte(keyBase)
	keyLength := byte(len(cryptKey))
	encrypted, err := cipher.Encrypt(out, securitySecret, cryptKey)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}

	// ヘッダー部(圧縮判定+暗号化キー長+暗号化キー)のバイト長を保存
	headerLength := make([]byte, headerLengthBytes)
	binary.LittleEndian.PutUint16(headerLength, uint16(metadataBytes+keyLength))

	totalBytes := headerLengthBytes + metadataBytes + int(keyLength) + len(encrypted)
	result := make([]byte, 0, totalBytes)
	result = append(result, headerLength...)
	result = append(result, byte(compress), keyLength)
	result = append(result, cryptKey...)
	result = append(result, encrypted...)

	return result, nil
}

func (Codec) Unmarshal(data []byte, v interface{}) (err error) {
	defer func() {
		// Panicが発生した場合でもプロセスを落とさずエラーとして処理する
		if r := recover(); r != nil {
			err = cerrors.Newf(cerrors.Internal, "panic has occurred in unmarshaling. %v", r)
		}
	}()

	// メッセージデータが空だった場合はunmarshal処理をしない.(インスタンスの初期値を利用する)
	if len(data) == 0 {
		return nil
	}

	// バイト数がヘッダー長2バイト+圧縮判定1バイト+暗号化key長1バイト=4バイト以下であれば不正なデータ
	if len(data) <= headerLengthBytes+metadataBytes {
		return cerrors.Newf(cerrors.Internal, "data was too short")
	}
	headerLength := int(binary.LittleEndian.Uint16(data[0:headerLengthBytes]))
	messageData := data[headerLengthBytes:]

	isCompressed := int(messageData[0]) == 1
	keyLength := int(messageData[1])
	if keyLength == 0 {
		return cerrors.Newf(cerrors.Internal, "crypt key not found")
	}

	// メッセージデータのバイト長が圧縮判定1バイト+暗号化key長1バイト+暗号化key長バイト以下であればBody部分のない不正なデータ
	if len(messageData) <= headerLength {
		return cerrors.Newf(cerrors.Internal, "message data was too short")
	}

	cryptKey := messageData[metadataBytes:headerLength]
	encrypted := messageData[headerLength:]

	// 暗号化されたMessageデータの復号
	decrypted, err := cipher.Decrypt(encrypted, securitySecret, cryptKey)
	if err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}

	// 圧縮されたデータの解凍
	if isCompressed {
		reader, err := gzip.NewReader(bytes.NewBuffer(decrypted))
		if err != nil {
			return cerrors.Wrap(err, cerrors.Internal)
		}
		defer func() {
			if err := reader.Close(); err != nil {
				// エラーログを吐いて処理は進める
				app.GetLogger().Error(context.Background(), err.Error())
			}
		}()

		buffer := new(bytes.Buffer)
		if _, err := buffer.ReadFrom(reader); err != nil {
			return cerrors.Wrap(err, cerrors.Internal)
		}
		decrypted = buffer.Bytes()
	}

	// proto unmarshal
	if err := unmarshal(decrypted, v); err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

// Name returns the name of the codec implementation. The returned string
// will be used as part of content type in transmission.  The result must be
// static; the result cannot change between calls.
func (Codec) Name() string {
	return Name
}

func (Codec) String() string {
	return Name
}
