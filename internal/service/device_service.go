package service

import (
	"context"
	"errors"

	model "github.com/Kowari1/File-Handler/internal/domain"
	"github.com/Kowari1/File-Handler/internal/repository"
	"github.com/google/uuid"
)

type DeviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

type PageResult struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
	Data  []model.Device `json:"data"`
}

func (s *DeviceService) GetAll(
	ctx context.Context,
	page, limit int,
) (*PageResult, error) {

	if page < 1 {
		return nil, errors.New("page must be >= 1")
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	total, err := s.repo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	devices, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return &PageResult{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  devices,
	}, nil
}

func (s *DeviceService) GetByUnitGUID(
	ctx context.Context,
	guid uuid.UUID,
	page, limit int,
) (*PageResult, error) {

	if page < 1 {
		return nil, errors.New("page must be >= 1")
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	total, err := s.repo.CountByUnitGUID(ctx, guid)
	if err != nil {
		return nil, err
	}

	devices, err := s.repo.FindLimitedByUnitGUID(ctx, guid, limit, offset)
	if err != nil {
		return nil, err
	}

	return &PageResult{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  devices,
	}, nil
}
