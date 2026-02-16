package repository

import (
	"context"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/google/uuid"
)

type DeviceRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]model.Device, error)
	FindLimitedByUnitGUID(ctx context.Context, guid uuid.UUID, limit, offset int) ([]model.Device, error)
	CountAll(ctx context.Context) (int, error)
	CountByUnitGUID(ctx context.Context, guid uuid.UUID) (int, error)
	Save(ctx context.Context, device []*model.Device) error
	FindByUnitGUID(ctx context.Context, guid uuid.UUID) ([]*model.Device, error)
}

type ProcessedFileRepository interface {
	Create(ctx context.Context, fileName string) error
	UpdateStatus(ctx context.Context, fileName string, status string, errMsg *string) error
	Exists(ctx context.Context, fileName string) (bool, error)
}

type ParseErrorRepository interface {
	Save(ctx context.Context, fileName string, line int, message string) error
}
