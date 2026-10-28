package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/guv41k/todo_app/internal/core/domain"
	core_errors "github.com/guv41k/todo_app/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET
		full_name=$1,
		phone_number=$2,
		version=version+1
	WHERE id=$3 AND version=$4
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, user.ID, user.Version)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' concurrently accessed: %w", user.ID, core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan updated user: %w", err)
	}

	return modelToDomain(userModel), nil
}
