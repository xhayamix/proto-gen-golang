{{ template "autogen_comment" }}
package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/QualiArts/campus-server/pkg/domain/proto/client/api"
)

func registerHandler(ctx context.Context, mux *runtime.ServeMux, target string, opts []grpc.DialOption) error {
	endpoints := []func(ctx context.Context, mux *runtime.ServeMux, target string, opts []grpc.DialOption) error{
	{{- range . }}
		pb.Register{{ .PascalName }}HandlerFromEndpoint,
	{{- end }}
	}
	for _, v := range endpoints {
		if err := v(ctx, mux, target, opts); err != nil {
			return err
		}
	}
	return nil
}
