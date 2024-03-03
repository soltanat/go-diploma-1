package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soltanat/go-diploma-1/internal/entities"
	"github.com/soltanat/go-diploma-1/internal/usecases/storager"
)

type OrderStorage struct {
	conn *pgxpool.Pool
}

func NewOrderStorage(conn *pgxpool.Pool) storager.OrderStorager {
	return &OrderStorage{
		conn: conn,
	}
}

func (s *OrderStorage) Save(ctx context.Context, tx storager.Tx, order *entities.Order) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	_, err := conn.Exec(ctx, `INSERT INTO service_diploma_1.orders (number, status, whole, decimal, uploaded_at, user_id) VALUES ($1, $2, $3, $4, $5)`, order.Number, order.Accrual.Whole, order.Accrual.Decimal, order.UploadedAt, order.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.ExistOrderError{}
		}
		return entities.StorageError{Err: err}
	}
	return nil
}

func (s *OrderStorage) Get(ctx context.Context, tx storager.Tx, number entities.OrderNumber) (*entities.Order, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	row := conn.QueryRow(ctx, `SELECT number, status, whole, decimal, uploaded_at, user_id FROM service_diploma_1.orders WHERE number = $1`, number)
	var order *entities.Order
	if err := row.Scan(order.Number, order.Accrual.Whole, order.Accrual.Decimal, order.UploadedAt, order.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.NotFoundError{}
		}
		return nil, entities.StorageError{Err: err}
	}

	return order, nil
}

func (s *OrderStorage) List(ctx context.Context, tx storager.Tx, userID *entities.Login, status *[]entities.OrderStatus) ([]entities.Order, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	var rows pgx.Rows
	var err error
	if userID != nil && status != nil {
		rows, err = conn.Query(ctx, `SELECT number, status, whole, decimal, uploaded_at, user_id FROM service_diploma_1.orders WHERE user_id = $1`, userID)
		if err != nil {
			return nil, entities.StorageError{Err: err}
		}
	} else if userID != nil {
		rows, err = conn.Query(ctx, `SELECT number, status, whole, decimal, uploaded_at, user_id FROM service_diploma_1.orders WHERE user_id = $1`, userID)
		if err != nil {
			return nil, entities.StorageError{Err: err}
		}
	} else if status != nil {
		rows, err = conn.Query(ctx, `SELECT number, status, whole, decimal, uploaded_at, user_id FROM service_diploma_1.orders WHERE status = ANY($1)`, status)
		if err != nil {
			return nil, entities.StorageError{Err: err}
		}
	}

	var orders []entities.Order
	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(&order.Number, &order.Accrual.Whole, &order.Accrual.Decimal, &order.UploadedAt, &order.UserID); err != nil {
			return nil, entities.StorageError{Err: err}
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderStorage) Update(ctx context.Context, tx storager.Tx, order *entities.Order) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	_, err := conn.Exec(ctx, `UPDATE service_diploma_1.orders SET status = $1, whole = $2, decimal = $3, uploaded_at = $4 WHERE number = $5`, order.Status, order.Accrual.Whole, order.Accrual.Decimal, order.UploadedAt, order.Number)
	if err != nil {
		return entities.StorageError{Err: err}
	}

	return nil
}

func (s *OrderStorage) Tx(ctx context.Context) storager.Tx {
	return &Tx{
		conn: s.conn,
		tx:   nil,
	}
}
