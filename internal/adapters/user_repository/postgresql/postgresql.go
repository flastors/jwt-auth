package postgresql

import (
	"context"
	"strings"

	"github.com/flastors/jwt-auth-golang/internal/core/user"
	postgresql "github.com/flastors/jwt-auth-golang/pkg/client/postgres"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
)

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

type userRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewUserRepository(client postgresql.Client, logger *logging.Logger) user.UserRepository {
	return &userRepository{
		client: client,
		logger: logger,
	}
}

func (r *userRepository) GetByID(ctx context.Context, userId string) (*user.User, error) {
	q := `
		SELECT id, email, refresh_token
		FROM public.user
		WHERE id = $1
	`
	u := &user.User{}
	r.logger.Trace(formatQuery(q))
	err := r.client.QueryRow(ctx, q, userId).Scan(&u.ID, &u.Email, &u.RefreshToken)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) Save(ctx context.Context, user *user.User) error {
	user2, err := r.GetByID(ctx, user.ID)
	if err != nil || user2.ID == "" {
		email := "john@doe.com"
		q := `
			INSERT INTO public.user (id, email, refresh_token) 
			VALUES ($1, $2, $3)
		`
		r.logger.Trace(formatQuery(q))
		_, err := r.client.Exec(ctx, q, user.ID, email, user.RefreshToken)
		if err != nil {
			return err
		}
	} else {
		err := r.UpdateRefreshToken(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) UpdateRefreshToken(ctx context.Context, user *user.User) error {
	q := `
		UPDATE public.user 
		SET refresh_token = $2 
		WHERE id = $1
	`
	r.logger.Trace(formatQuery(q))
	_, err := r.client.Exec(ctx, q, user.ID, user.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}
