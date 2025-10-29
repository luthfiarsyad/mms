package mysql

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/luthfiarsyad/mms/internal/domain/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }
func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	q := `INSERT INTO users (name, email, password, created_at) VALUES
(?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, q, u.Name, u.Email, u.Password,
		u.CreatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT id, name, email, password, created_at FROM users WHERE email
= ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &u, nil
}
func (r *UserRepo) FindByID(ctx context.Context, id int64) (*domain.User,
	error) {
	q := `SELECT id, name, email, password, created_at FROM users WHERE id = ?
LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
