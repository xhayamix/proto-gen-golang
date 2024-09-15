package genserverapi

import (
	"crypto/tls"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/port/api"
	"github.com/xhayamix/proto-gen-golang/pkg/grpc/codec/protoenc"
)

type client struct {
	host     string
	grpcOpts []grpc.DialOption
	conn     *grpc.ClientConn
}

func New(host, secret string, insec bool) (api.API, error) {
	// codec登録
	protoenc.Register([]byte(secret))
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			// NOTE: api側でgrpc.CustomCodecをしている間はForceCodecが必要
			grpc.ForceCodec(protoenc.Codec{}),
			grpc.CallContentSubtype(protoenc.Name),
		),
	}

	var cred credentials.TransportCredentials
	if insec {
		cred = insecure.NewCredentials()
	} else {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		cred = credentials.NewTLS(&tls.Config{
			RootCAs:    systemRoots,
			MinVersion: tls.VersionTLS12,
		})
	}
	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(cred))
	c := &client{host: host, grpcOpts: grpcOpts}
	return c, nil
}

func (c *client) Close() error {
	if c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return cerrors.Wrap(err, cerrors.Internal)
	}
	return nil
}

func (c *client) getConn() (*grpc.ClientConn, error) {
	// 生きているコネクションがなければ再接続する
	if c.conn == nil || c.conn.GetState() == connectivity.Shutdown {
		conn, err := grpc.Dial(c.host, c.grpcOpts...)
		if err != nil {
			return nil, cerrors.Wrap(err, cerrors.Internal)
		}
		c.conn = conn
	}

	return c.conn, nil
}
