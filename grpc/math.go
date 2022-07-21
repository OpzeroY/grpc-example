package main

import (
	"context"
	grpc_math "test/grpc/protocol/math"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func newGRPCServer(mux *runtime.ServeMux, cc *grpc.ClientConn) *grpc.Server {
	s := grpc.NewServer()

	// register to grpc
	grpc_math.RegisterCatServer(s, Math{})
	if mux != nil && cc != nil {
		// register to gateway
		grpc_math.RegisterCatHandler(context.Background(), mux, cc)
	}
	return s
}

type Math struct {
	grpc_math.UnimplementedCatServer
}

func (Math) Sum(ctx context.Context, req *grpc_math.SumRequest) (*grpc_math.SumResponse, error) {
	var sum int32 = 0
	for _, val := range req.Vals {
		sum += val
	}
	return &grpc_math.SumResponse{
		Val: sum,
	}, nil
}
func (c Math) Version(context.Context, *grpc_math.VersionRequest) (*grpc_math.VersionResponse, error) {
	return &grpc_math.VersionResponse{
		Val: `v1.0.0`,
	}, nil
}
