// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package registry

import (
	"google.golang.org/grpc"

	pb "github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api"
)

func registerServer(s *grpc.Server) func(
	userServer pb.UserServer,
) {
	return func(
		userServer pb.UserServer,
	) {
		pb.RegisterUserServer(s, userServer)
	}
}
