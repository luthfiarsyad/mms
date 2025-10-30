package usecase

import (
	"context"

	"github.com/luthfiarsyad/mms/internal/domain/transaction"
	"github.com/luthfiarsyad/mms/internal/infrastructure/logger"
)

type TransactionUsecase struct {
	txService *transaction.Service
}

func NewTransactionUsecase(txService *transaction.Service) *TransactionUsecase {
	logger.L.Debug().Msg("TransactionUsecase: initialized")
	return &TransactionUsecase{txService: txService}
}

func (u *TransactionUsecase) CreateTransaction(ctx context.Context, userID int64, amount float64, description, txType string) (*transaction.Transaction, error) {
	logger.L.Info().
		Int64("user_id", userID).
		Float64("amount", amount).
		Str("type", txType).
		Msg("TransactionUsecase.CreateTransaction: creating transaction")

	t := &transaction.Transaction{
		UserID:      userID,
		Amount:      amount,
		Description: description,
		Type:        txType,
	}

	err := u.txService.Create(ctx, t)
	if err != nil {
		logger.L.Error().
			Err(err).
			Int64("user_id", userID).
			Float64("amount", amount).
			Str("type", txType).
			Msg("TransactionUsecase.CreateTransaction: failed to create transaction")
		return nil, err
	}

	logger.L.Info().
		Int64("transaction_id", t.ID).
		Int64("user_id", userID).
		Float64("amount", amount).
		Str("type", txType).
		Msg("TransactionUsecase.CreateTransaction: transaction created successfully")

	return t, nil
}

func (u *TransactionUsecase) GetTransactionByID(ctx context.Context, id int64) (*transaction.Transaction, error) {
	logger.L.Info().
		Int64("transaction_id", id).
		Msg("TransactionUsecase.GetTransactionByID: fetching transaction")

	t, err := u.txService.GetByID(ctx, id)
	if err != nil {
		logger.L.Error().
			Err(err).
			Int64("transaction_id", id).
			Msg("TransactionUsecase.GetTransactionByID: failed to fetch transaction")
		return nil, err
	}

	logger.L.Info().
		Int64("transaction_id", id).
		Msg("TransactionUsecase.GetTransactionByID: transaction fetched successfully")

	return t, nil
}

func (u *TransactionUsecase) GetUserTransactions(ctx context.Context, userID int64) ([]*transaction.Transaction, error) {
	logger.L.Info().
		Int64("user_id", userID).
		Msg("TransactionUsecase.GetUserTransactions: fetching user transactions")

	transactions, err := u.txService.GetByUserID(ctx, userID)
	if err != nil {
		logger.L.Error().
			Err(err).
			Int64("user_id", userID).
			Msg("TransactionUsecase.GetUserTransactions: failed to fetch user transactions")
		return nil, err
	}

	logger.L.Info().
		Int64("user_id", userID).
		Int("count", len(transactions)).
		Msg("TransactionUsecase.GetUserTransactions: transactions fetched successfully")

	return transactions, nil
}

func (u *TransactionUsecase) UpdateTransaction(ctx context.Context, id int64, amount float64, description, txType string) (*transaction.Transaction, error) {
	logger.L.Info().
		Int64("transaction_id", id).
		Float64("amount", amount).
		Str("type", txType).
		Msg("TransactionUsecase.UpdateTransaction: updating transaction")

	t := &transaction.Transaction{
		ID:          id,
		Amount:      amount,
		Description: description,
		Type:        txType,
	}

	err := u.txService.Update(ctx, t)
	if err != nil {
		logger.L.Error().
			Err(err).
			Int64("transaction_id", id).
			Msg("TransactionUsecase.UpdateTransaction: failed to update transaction")
		return nil, err
	}

	logger.L.Info().
		Int64("transaction_id", id).
		Msg("TransactionUsecase.UpdateTransaction: transaction updated successfully")

	return t, nil
}

func (u *TransactionUsecase) DeleteTransaction(ctx context.Context, id int64) error {
	logger.L.Info().
		Int64("transaction_id", id).
		Msg("TransactionUsecase.DeleteTransaction: deleting transaction")

	err := u.txService.Delete(ctx, id)
	if err != nil {
		logger.L.Error().
			Err(err).
			Int64("transaction_id", id).
			Msg("TransactionUsecase.DeleteTransaction: failed to delete transaction")
		return err
	}

	logger.L.Info().
		Int64("transaction_id", id).
		Msg("TransactionUsecase.DeleteTransaction: transaction deleted successfully")

	return nil
}
