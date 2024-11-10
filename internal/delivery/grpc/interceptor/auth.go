package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/Karzoug/meower-user-service/pkg/auth"
)

const userKey string = "x-user"

func Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, auth.NewAuthNError()
		}
		key, found := md[userKey]
		if !found || len(key) == 0 {
			return nil, auth.NewAuthNError()
		}

		ctx = auth.WithUserID(ctx, key[0])

		return handler(ctx, req)
	}
}
