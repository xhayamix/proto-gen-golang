// Code generated by protoc-gen-all. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package genserverapi

import (
	"context"

	"github.com/xhayamix/proto-gen-golang/pkg/cerrors"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/proto/client/api"
)

func (c *client) UserServiceGetProfile(ctx context.Context, req *api.UserServiceGetProfileRequest) (*api.UserServiceGetProfileResponse, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, cerrors.Stack(err)
	}
	cli := api.NewUserServiceClient(conn)
	result, err := cli.GetProfile(ctx, req)
	if err != nil {
		return nil, cerrors.Wrap(err, cerrors.Internal)
	}
	return result, nil
}
