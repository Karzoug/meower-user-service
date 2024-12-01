package converter

import (
	gen "github.com/Karzoug/meower-user-service/internal/delivery/grpc/gen/user/v1"
	"github.com/Karzoug/meower-user-service/internal/user/entity"
)

func ToProtoUser(u entity.User) *gen.User {
	return &gen.User{
		Id:         u.ID.String(),
		Username:   u.Username,
		Name:       u.Name,
		ImageUrl:   u.ImageURL,
		StatusText: u.StatusText,
	}
}

func ToProtoUserShortProjection(u entity.UserShortProjection) *gen.UserShortProjection {
	return &gen.UserShortProjection{
		Id:         u.ID.String(),
		Username:   u.Username,
		Name:       u.Name,
		ImageUrl:   u.ImageURL,
		StatusText: u.StatusText,
	}
}

func ToProtoUserShortProjections(users []entity.UserShortProjection) []*gen.UserShortProjection {
	projections := make([]*gen.UserShortProjection, len(users))
	for i := range users {
		projections[i] = ToProtoUserShortProjection(users[i])
	}

	return projections
}
