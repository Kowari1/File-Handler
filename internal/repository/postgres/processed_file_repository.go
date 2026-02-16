package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcessedFileRepository struct {
	pool *pgxpool.Pool
}

func NewProcessedFileRepository(pool *pgxpool.Pool) *ProcessedFileRepository {
	return &ProcessedFileRepository{
		pool: pool,
	}
}

func (r *ProcessedFileRepository) Create(
	ctx context.Context,
	fileName string,
) error {

	query := `
		INSERT INTO processed_file (file_name, status)
		VALUES ($1, 'processing')
	`

	_, err := r.pool.Exec(ctx, query, fileName)
	return err
}

func (r *ProcessedFileRepository) UpdateStatus(
	ctx context.Context,
	fileName string,
	status string,
	errMsg *string,
) error {

	query := `
		UPDATE processed_file
		SET status = $1,
		    error_message = $2
		WHERE file_name = $3
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		status,
		errMsg,
		fileName,
	)

	return err
}

func (r *ProcessedFileRepository) Exists(
	ctx context.Context,
	fileName string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1 FROM processed_file WHERE file_name = $1
		)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, fileName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
