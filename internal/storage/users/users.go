package users

import (
	"context"
	stdsql "database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/getbread/gokit/storage/sql"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID        uuid.UUID `db:"id"`
	UserName  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateUser struct {
	ID        uuid.UUID `db:"id"`
	UserName  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
}

type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	CreateUser(ctx context.Context, c CreateUser) (User, error)
}

type store struct {
	db sql.Adapter
}

const getUserByIDQuery = `
	SELECT id, username, first_name, last_name, created_at, updated_at FROM users WHERE id = $1;
`

func (s store) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	if err := s.db.GetContext(ctx, &user, getUserByIDQuery, id); err != nil {
		if errors.Is(err, stdsql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

const insertUserQuery = `
	INSERT INTO users (id, username, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING  created_at, updated_at;
`

func (s store) CreateUser(ctx context.Context, c CreateUser) (User, error) {
	var createdAt, updatedAt time.Time
	if err := s.db.QueryRowContext(ctx, insertUserQuery, c.ID, c.UserName, c.FirstName, c.LastName).Scan(&createdAt, &updatedAt); err != nil {
		return User{}, err
	}

	return User{
		ID:        c.ID,
		UserName:  c.UserName,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func NewStore(db sql.Adapter) Repository {
	return &store{db: db}
}
