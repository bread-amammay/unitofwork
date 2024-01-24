package storage

import (
	"context"

	"github.com/bread-amammay/unitofwork/internal/storage/blogs"
	"github.com/bread-amammay/unitofwork/internal/storage/users"
	"github.com/jmoiron/sqlx"

	"github.com/getbread/gokit/db"
)

type UnitOfWorkBlock func(context.Context, UnitOfWorkStore) error

type UnitOfWorkStore interface {
	Blogs() blogs.Repository
	Users() users.Repository
}

type uowStore struct {
	blogs blogs.Repository
	users users.Repository
}

func (u uowStore) Blogs() blogs.Repository {
	return u.blogs
}

func (u uowStore) Users() users.Repository {
	return u.users
}

type UnitOfWork interface {
	Do(context.Context, UnitOfWorkBlock) error
}

type unitOfWork struct {
	db   *sqlx.DB
	txer db.Transactioner
}

func New(d *sqlx.DB) UnitOfWork {
	return &unitOfWork{db: d, txer: db.NewTransactioner(d)}
}

// Do execute the given UnitOfWorkBlock.
func (s *unitOfWork) Do(ctx context.Context, fn UnitOfWorkBlock) error {
	return s.txer.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		newStore := &uowStore{
			blogs: blogs.NewStore(tx),
			users: users.NewStore(tx),
		}
		return fn(ctx, newStore)
	})

}
