package pg

import (
	"github.com/Karzoug/meower-common-go/postgresql"
)

type repo struct {
	db postgresql.DB
}

func NewUserRepo(db postgresql.DB) repo {
	return repo{db: db}
}
