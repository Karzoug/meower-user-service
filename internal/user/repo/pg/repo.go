package pg

import (
	"github.com/Karzoug/meower-user-service/pkg/postgresql"
)

const tableName = "users"

type repo struct {
	db postgresql.DB
}

func NewUserRepo(db postgresql.DB) repo {
	return repo{db: db}
}
