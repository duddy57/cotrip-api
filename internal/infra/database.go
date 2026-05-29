package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context, user, pass, host, port, db string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		user, pass, host, port, db,
	))
	if err != nil {
		return nil, err
	}
	return pool, nil
}
