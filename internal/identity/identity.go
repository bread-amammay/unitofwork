package identity

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Identity struct {
	UserID    uuid.UUID
	UserName  string
	FirstName string
	LastName  string
}

var contextKey = struct{}{}

func OnContext(ctx context.Context) Identity {
	return ctx.Value(contextKey).(Identity)
}

func WithIdentity(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, contextKey, identity)
}

func ConnectInterceptor(z zerolog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {

			userID := request.Header().Get("X-User-Id")
			userName := request.Header().Get("X-User-Name")
			firstName := request.Header().Get("X-First-Name")
			lastName := request.Header().Get("X-Last-Name")

			if userID == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing user id"))
			}
			id, err := uuid.Parse(userID)
			if err != nil {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}

			if userName == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing user name"))
			}

			if firstName == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing first name"))
			}

			if lastName == "" {
				return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("missing last name"))
			}

			identity := Identity{
				UserID:    id,
				UserName:  userName,
				FirstName: firstName,
				LastName:  lastName,
			}

			logger := z.With().Str("user_id", identity.UserID.String()).Str("user_name", identity.UserName).Str("first_name", identity.FirstName).Str("last_name", identity.LastName).Logger()
			ctx = logger.WithContext(ctx)
			ctx = WithIdentity(ctx, identity)
			return next(ctx, request)
		})

	}
}
