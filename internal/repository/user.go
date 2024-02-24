package repository

import (
	"auth/internal/domain"
	"auth/internal/repository/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepos struct {
}

func NewUserRepos() User {
	return &UserRepos{}
}

func (m *UserRepos) AddUser(ctx context.Context, tx pgx.Tx, email string, role domain.Role, passHash []byte) (int, error) {
	row := tx.QueryRow(ctx, `INSERT INTO "user" (email, pass_hash, role) VALUES ($1, $2, $3) RETURNING id`,
		email, passHash, role)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == postgres.ErrUniqueViolation {
			return 0, UserAlreadyExist
		}
		return 0, fmt.Errorf("AddUser/Scan: %w", err)
	}
	return id, nil
}

func (m *UserRepos) UserById(ctx context.Context, tx pgx.Tx, id int) (*domain.UserFromDB, error) {
	row := tx.QueryRow(ctx, `SELECT email, pass_hash, role FROM "user" WHERE id = $1`, id)

	var user domain.UserFromDB
	err := row.Scan(
		&user.Email,
		&user.HashPassword,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotExist
		}
		return nil, fmt.Errorf("UserById/Scan: %w", err)
	}
	user.ID = id
	return &user, nil
}

func (m *UserRepos) UserByEmail(ctx context.Context, tx pgx.Tx, email string) (*domain.UserFromDB, error) {
	row := tx.QueryRow(ctx, `SELECT id, pass_hash, role FROM "user" WHERE email = $1`, email)

	var user domain.UserFromDB
	err := row.Scan(
		&user.ID,
		&user.HashPassword,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotExist
		}
		return nil, fmt.Errorf("UserById/Scan: %w", err)
	}
	user.Email = email
	return &user, nil
}
