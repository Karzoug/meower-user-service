package health

import (
	"context"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func RegisterService() func(grpcServer *grpc.Server) {
	hdl := handlers{}

	return func(grpcServer *grpc.Server) {
		health.RegisterHealthServer(grpcServer, hdl)
	}
}

type handlers struct {
	health.UnimplementedHealthServer
}

func (h handlers) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

func (h handlers) Watch(req *health.HealthCheckRequest, ss grpc.ServerStreamingServer[health.HealthCheckResponse]) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}
