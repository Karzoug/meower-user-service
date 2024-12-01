package user

import (
	"context"

	"google.golang.org/grpc"

	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Karzoug/meower-user-service/internal/delivery/grpc/converter"
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
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	id, err := xid.FromString(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id: "+req.Id)
	}

	user, err := h.userService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return converter.ToProtoUser(user), nil
}

func (h handlers) GetShortProjection(ctx context.Context, req *gen.GetShortProjectionRequest) (*gen.UserShortProjection, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if reqID := req.GetId(); reqID != "" {
		id, err := xid.FromString(reqID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid id: "+reqID)
		}

		projection, err := h.userService.GetShortProjection(ctx, id)
		if err != nil {
			return nil, err
		}

		return converter.ToProtoUserShortProjection(projection), nil
	}

	if reqUsername := req.GetUsername(); reqUsername != "" {
		projection, err := h.userService.GetShortProjectionByUsername(ctx, reqUsername)
		if err != nil {
			return nil, err
		}

		return converter.ToProtoUserShortProjection(projection), nil
	}

	return nil, status.Error(codes.InvalidArgument, "empty id or username")
}

func (h handlers) BatchGetShortProjections(ctx context.Context, req *gen.BatchGetShortProjectionsRequest) (*gen.BatchGetShortProjectionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ids := make([]xid.ID, len(req.Ids))
	var err error
	for i, id := range req.Ids {
		ids[i], err = xid.FromString(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid id: "+id)
		}
	}

	users, err := h.userService.BatchGetShortProjections(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &gen.BatchGetShortProjectionsResponse{
		Users: converter.ToProtoUserShortProjections(users),
	}, nil
}
