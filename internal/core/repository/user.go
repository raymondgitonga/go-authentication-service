package repository

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) AddUser(name string, secret []byte) error {
	query := `INSERT INTO service_user (name, secret) VALUES ($1,$2);`

	_, err := r.db.Exec(query, name, secret)

	if err != nil {
		return fmt.Errorf("error adding value to db %w", err)
	}
	return nil
}
