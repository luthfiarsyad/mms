package mysql

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/luthfiarsyad/mms/internal/domain/transaction"
)

type TxRepo struct {
	db *sql.DB
}

func NewTxRepo(db *sql.DB) *TxRepo {
	return &TxRepo{db: db}
}

func (r *TxRepo) Create(ctx context.Context, t *domain.Transaction) error {
	q := `INSERT INTO transactions (user_id, amount, description, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, q, t.UserID, t.Amount, t.Description, t.Type, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (r *TxRepo) FindByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	q := `SELECT id, user_id, amount, description, type, created_at, updated_at FROM transactions WHERE id = ? LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, id)
	var t domain.Transaction
	if err := row.Scan(&t.ID, &t.UserID, &t.Amount, &t.Description, &t.Type, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &t, nil
}

func (r *TxRepo) FindByUserID(ctx context.Context, userID int64) ([]*domain.Transaction, error) {
	q := `SELECT id, user_id, amount, description, type, created_at, updated_at FROM transactions WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Description, &t.Type, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TxRepo) Update(ctx context.Context, t *domain.Transaction) error {
	q := `UPDATE transactions SET amount = ?, description = ?, type = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, t.Amount, t.Description, t.Type, t.UpdatedAt, t.ID)
	return err
}

func (r *TxRepo) Delete(ctx context.Context, id int64) error {
	q := `DELETE FROM transactions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}