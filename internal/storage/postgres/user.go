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

type UserStorage struct {
	conn *pgxpool.Pool
}

func NewUserStorage(conn *pgxpool.Pool) storager.UserStorager {
	return &UserStorage{
		conn: conn,
	}
}

func (s *UserStorage) Save(ctx context.Context, tx storager.Tx, user *entities.User) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	_, err := conn.Exec(
		ctx,
		`INSERT INTO users (login, password, whole, decimal) VALUES ($1, $2, $3, $4)`,
		user.Login, user.Password, user.Balance.Whole, user.Balance.Decimal,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entities.ExistUserError{}
		}
		return entities.StorageError{Err: err}
	}

	return nil
}

func (s *UserStorage) Get(ctx context.Context, tx storager.Tx, login entities.Login) (*entities.User, error) {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	row := conn.QueryRow(
		ctx,
		`SELECT login, password, whole, decimal FROM users WHERE login = $1`, login,
	)
	var user entities.User
	if err := row.Scan(&user.Login, &user.Password, &user.Balance.Whole, &user.Balance.Decimal); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entities.NotFoundError{}
		}
		return nil, entities.StorageError{Err: err}
	}

	return &user, nil
}

func (s *UserStorage) Update(ctx context.Context, tx storager.Tx, user *entities.User) error {
	conn := s.conn
	if tx != nil {
		conn = tx.(*Tx).conn
	}

	_, err := conn.Exec(ctx, `UPDATE users SET whole = $1, decimal = $2 WHERE login = $3`, user.Balance.Whole, user.Balance.Decimal, user.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.NotFoundError{}
		}
		return entities.StorageError{Err: err}
	}

	return nil
}

func (s *UserStorage) Tx(ctx context.Context) storager.Tx {
	return &Tx{
		conn: s.conn,
		tx:   nil,
	}
}
