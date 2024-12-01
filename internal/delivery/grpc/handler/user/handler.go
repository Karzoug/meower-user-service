package user

import (
	"context"

	"google.golang.org/grpc"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/grpc/gen/user/v1"
	"github.com/Karzoug/meower-user-service/internal/user/service"
)

func RegisterService(us service.UserService) func(grpcServer *grpc.Server) {
	hdl := handlers{
		userService: us,
	}
	return func(grpcServer *grpc.Server) {
		gen.RegisterUserServiceServer(grpcServer, hdl)
	}
}

type handlers struct {
	gen.UnimplementedUserServiceServer
	userService service.UserService
}

func (h handlers) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.User, error) {
	panic("not implemented")
}

func (h handlers) GetShortProjection(ctx context.Context, req *gen.GetShortProjectionRequest) (*gen.UserShortProjection, error) {
	panic("not implemented")
}

func (h handlers) BatchGetShortProjections(ctx context.Context, req *gen.BatchGetShortProjectionsRequest) (*gen.BatchGetShortProjectionsResponse, error) {
	panic("not implemented")
}
