package repositories

import (
	"context"
	"database/sql"
	"spending/utils"

	"go.opentelemetry.io/otel/trace"
)

type DbTx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
}

func Query[T any](span trace.Span, query func() (*sql.Rows, error), read func(*sql.Rows) *T) (*T, error) {
	rows, err := query()
	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	result := read(rows)

	return result, err
}

func QueryList[T any](span trace.Span, query func() (*sql.Rows, error), read func(*sql.Rows) *T) ([]*T, error) {
	rows, err := query()

	if err != nil {
		utils.TraceError(span, err)
		return nil, err
	}

	defer rows.Close()

	var results []*T = make([]*T, 0)

	for rows.Next() {
		result := read(rows)
		if result != nil {
			results = append(results, result)
		}
	}

	return results, err
}
