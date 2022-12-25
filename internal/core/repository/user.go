package repository

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) AddUser(name string, secret []byte) error {
	query := `INSERT INTO service_user (name, secret) VALUES ($1,$2);`

	_, err := r.db.Exec(query, name, secret)

	if err != nil {
		r.logger.Error("error at AddUser", zap.String("error", err.Error()))
		return fmt.Errorf("error adding user")
	}
	return nil
}

func (r *UserRepository) GetUser(name string) (string, error) {
	var secret string
	query := `SELECT secret from service_user WHERE name=$1;`

	row := r.db.QueryRow(query, name)

	err := row.Scan(&secret)
	if err != nil {
		r.logger.Error("error at  GetUser", zap.String("error", err.Error()))
		return "", fmt.Errorf("error getting user")
	}
	return secret, nil
}
