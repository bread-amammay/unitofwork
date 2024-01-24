package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bread-amammay/unitofwork/internal/storage"
	"github.com/bread-amammay/unitofwork/internal/storage/blogs"
	"github.com/bread-amammay/unitofwork/internal/storage/users"
	"github.com/google/uuid"
)

type Controller struct {
	uow storage.UnitOfWork
}

func NewController(uow storage.UnitOfWork) Controller {
	return Controller{uow: uow}
}

type CreateBlog struct {
	Title string
	Body  string
	User  User
}
type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlogPost struct {
	ID        uuid.UUID
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}

func (c Controller) CreateBlog(ctx context.Context, create CreateBlog) (BlogPost, error) {

	var bp BlogPost

	err := c.uow.Do(ctx, func(ctx context.Context, store storage.UnitOfWorkStore) error {
		user, err := store.Users().GetUserByID(ctx, create.User.ID)
		if errors.Is(err, users.ErrUserNotFound) {
			user, err = store.Users().CreateUser(ctx, users.CreateUser{
				ID:        create.User.ID,
				UserName:  create.User.UserName,
				FirstName: create.User.FirstName,
				LastName:  create.User.LastName,
			})
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		insertBlog, err := store.Blogs().InsertBlog(ctx, blogs.CreateBlog{
			Title:  create.Title,
			Body:   create.Body,
			UserID: create.User.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to insert blog: %w", err)
		}

		bp = BlogPost{
			ID:        insertBlog.ID,
			Title:     insertBlog.Title,
			Body:      insertBlog.Body,
			CreatedAt: insertBlog.CreatedAt,
			UpdatedAt: insertBlog.UpdatedAt,
			User: User{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				UserName:  user.UserName,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		}

		return nil
	})
	if err != nil {
		return BlogPost{}, fmt.Errorf("failed to create blog: %w", err)
	}

	return bp, nil
}
