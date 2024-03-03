package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type WithdrawalStorage struct {
	conn *pgxpool.Pool
}

func NewWithdrawalStorage(conn *pgxpool.Pool) storager.WithdrawalStorager {
	return &WithdrawalStorage{
		conn: conn,
	}
}

func (s *WithdrawalStorage) Save(ctx context.Context, tx storager.Tx, withdraw *entities.Withdrawal) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	_, err := conn.Exec(
		ctx,
		`INSERT INTO service_diploma_1.withdrawals (number, whole, decimal, processed_at, user_id) VALUES ($1, $2, $3, $4, $5)`,
		withdraw.OrderNumber, withdraw.Sum.Whole, withdraw.Sum.Decimal, withdraw.ProcessedAt, withdraw.UserID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.ExistWithdrawalError{}
		}
		return entities.StorageError{Err: err}
	}

	return nil
}

func (s *WithdrawalStorage) List(ctx context.Context, tx storager.Tx, userID entities.Login) ([]entities.Withdrawal, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	rows, err := conn.Query(ctx, `SELECT number, whole, decimal, processed_at, user_id FROM service_diploma_1.withdrawals WHERE user_id = $1`, userID)
	if err != nil {
		return nil, entities.StorageError{Err: err}
	}

	var withdrawals []entities.Withdrawal
	for rows.Next() {
		var withdrawal entities.Withdrawal
		if err := rows.Scan(&withdrawal.OrderNumber, &withdrawal.Sum.Whole, &withdrawal.Sum.Decimal, &withdrawal.ProcessedAt, &withdrawal.UserID); err != nil {
			return nil, entities.StorageError{Err: err}
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil

}

func (s *WithdrawalStorage) Count(ctx context.Context, tx storager.Tx, userID entities.Login) (int, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	var count int
	err := conn.QueryRow(ctx, `SELECT COUNT(*) FROM service_diploma_1.withdrawals WHERE user_id = $1`, userID).Scan(&count)
	if err != nil {
		return 0, entities.StorageError{Err: err}
	}

	return count, nil
}

func (s *WithdrawalStorage) Tx(ctx context.Context) storager.Tx {
	return &Tx{
		conn: s.conn,
		tx:   nil,
	}
}
