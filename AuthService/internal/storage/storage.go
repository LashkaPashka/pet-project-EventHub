package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LashkaPashka/EventHub/AuthService/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewStorage(pool *pgxpool.Pool, logger *slog.Logger) *Storage {
	return &Storage{
		pool:   pool,
		logger: logger,
	}
}

func (s *Storage) FindByEmail(email string) *model.User {
	const op = "AuthService.internal.storage.FindByEmail"

	query := `SELECT email, password, username FROM users WHERE email = @email`

	args := pgx.NamedArgs{
		"email": email,
	}

	var user model.User
	err := s.pool.QueryRow(context.Background(), query, args).Scan(&user.Email, &user.Password, &user.Username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return  nil
		}
		s.logger.Error("Error query get parameters", slog.String("Error", op), slog.String("Error", err.Error()))
		return nil
	}

	return &user
}

func (s *Storage) Create(user *model.User) int {
	const op = "AuthService.internal.storage.Create"

	query := `INSERT INTO users (email, password, username) VALUES(@email, @password, @username)`

	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
		"username": user.Username,
	}

	_, err := s.pool.Exec(context.Background(), query, args)
	if err != nil {
		s.logger.Error("Invalid query SQL", slog.String("Error: ", op))
		return 0
	}

	return 0
}
