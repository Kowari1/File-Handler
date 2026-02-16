package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ParseErrorRepository struct {
	pool *pgxpool.Pool
}

func NewParseErrorRepository(pool *pgxpool.Pool) *ParseErrorRepository {
	return &ParseErrorRepository{
		pool: pool,
	}
}

func (r *ParseErrorRepository) Save(
	ctx context.Context,
	fileName string,
	line int,
	message string,
) error {

	query := `
		INSERT INTO parse_error (
			file_name,
			line_number,
			error_message
		)
		VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		fileName,
		line,
		message,
	)

	return err
}
