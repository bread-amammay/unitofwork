package api

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/bread-amammay/unitofwork/gen/api/blogs/v1"
	"github.com/bread-amammay/unitofwork/gen/api/blogs/v1/blogsv1connect"
	"github.com/bread-amammay/unitofwork/internal/identity"
	"github.com/bread-amammay/unitofwork/internal/usecase"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	useCase usecase.Controller
}

func NewServer(useCase usecase.Controller) Server {
	return Server{useCase: useCase}
}

func (s Server) CreateBlog(ctx context.Context, c *connect.Request[v1.CreateBlogRequest]) (*connect.Response[v1.CreateBlogResponse], error) {
	id := identity.OnContext(ctx)
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("received creating blog request")
	blog, err := s.useCase.CreateBlog(ctx, usecase.CreateBlog{
		Title: c.Msg.Title,
		Body:  c.Msg.Content,
		User: usecase.User{
			ID:        id.UserID,
			UserName:  id.UserName,
			FirstName: id.FirstName,
			LastName:  id.LastName,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create blog")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logger.Info().Msg("created blog")

	return &connect.Response[v1.CreateBlogResponse]{
		Msg: &v1.CreateBlogResponse{
			Blog: &v1.Blog{
				Id: &v1.UUID{Value: []byte(blog.ID.String())},
				Author: &v1.Author{
					Id:        &v1.UUID{Value: []byte(blog.User.ID.String())},
					UserName:  blog.User.UserName,
					FirstName: blog.User.FirstName,
					LastName:  blog.User.LastName,
					CreatedAt: timestamppb.New(blog.User.CreatedAt),
					UpdatedAt: timestamppb.New(blog.User.UpdatedAt),
				},
				Title:     blog.Title,
				Content:   blog.Body,
				CreatedAt: timestamppb.New(blog.CreatedAt),
				UpdatedAt: timestamppb.New(blog.UpdatedAt),
			},
		},
	}, nil
}

var (
	_ blogsv1connect.BlogServiceHandler = (*Server)(nil)
)
