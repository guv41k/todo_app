package users_postgres_repository

import (
	"github.com/guv41k/todo_app/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func (m *UserModel) Scan(row pgx.Row) error {
	return row.Scan(
		&m.ID,
		&m.Version,
		&m.FullName,
		&m.PhoneNumber,
	)
}

func modelToDomain(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.FullName,
		model.PhoneNumber,
	)
}

func modelsToDomains(models []UserModel) []domain.User {
	domains := make([]domain.User, len(models))

	for i, model := range models {
		domains[i] = modelToDomain(model)
	}

	return domains
}
