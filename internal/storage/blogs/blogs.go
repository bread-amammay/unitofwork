package blogs

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/getbread/gokit/storage/sql"
)

type Blog struct {
	ID        uuid.UUID `db:"id"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UserID    uuid.UUID `db:"user_id"`
}

type CreateBlog struct {
	Title  string    `db:"title"`
	Body   string    `db:"body"`
	UserID uuid.UUID `db:"user_id"`
}

type Repository interface {
	InsertBlog(ctx context.Context, blog CreateBlog) (Blog, error)
	SelectBlogByID(ctx context.Context, id uuid.UUID) (Blog, error)
	SelectAllBlogs(ctx context.Context) ([]Blog, error)
	UpdateBlog(ctx context.Context, blog Blog) error
}

type store struct {
	db sql.Adapter
}

const insertBlogQuery = `
	INSERT INTO blog_posts (title, body, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
`

func (s store) InsertBlog(ctx context.Context, blog CreateBlog) (Blog, error) {
	var id uuid.UUID
	var createdAt, updatedAt time.Time
	if err := s.db.QueryRowContext(ctx, insertBlogQuery, blog.Title, blog.Body, blog.UserID).Scan(&id, &createdAt, &updatedAt); err != nil {
		return Blog{}, err
	}

	return Blog{
		ID:        id,
		Title:     blog.Title,
		Body:      blog.Body,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    blog.UserID,
	}, nil
}

const selectBlogByIDQuery = `
	SELECT id, title, body, created_at, updated_at, user_id FROM blog_posts WHERE id = $1;
`

func (s store) SelectBlogByID(ctx context.Context, id uuid.UUID) (Blog, error) {
	var blog Blog
	if err := s.db.GetContext(ctx, &blog, selectBlogByIDQuery, id); err != nil {
		return Blog{}, err
	}
	return blog, nil
}

const selectAllBlogsQuery = `
	SELECT id, title, body, created_at, updated_at, user_id FROM blog_posts;
`

func (s store) SelectAllBlogs(ctx context.Context) ([]Blog, error) {
	var blogs []Blog
	if err := s.db.SelectContext(ctx, &blogs, selectAllBlogsQuery); err != nil {
		return []Blog{}, err
	}
	return blogs, nil
}

func (s store) UpdateBlog(ctx context.Context, blog Blog) error {
	// TODO implement me
	panic("implement me")
}

func NewStore(db sql.Adapter) Repository {
	return &store{db}
}
