package service

import (
	"context"

	"gomora/module/record/domain/entity"
	"gomora/module/record/domain/repository"
)

// RecordQueryService handles the record query service logic
type RecordQueryService struct {
	repository.RecordQueryRepositoryInterface
}

// GetRecords retrieves all records
func (service *RecordQueryService) GetRecords(ctx context.Context, page *uint) ([]entity.Record, uint, error) {
	res, totalCount, err := service.RecordQueryRepositoryInterface.SelectRecords(page)
	if err != nil {
		return nil, 0, err
	}

	return res, totalCount, nil
}

// GetRecordByID retrieves the record provided by its id
func (service *RecordQueryService) GetRecordByID(ctx context.Context, ID string) (entity.Record, error) {
	res, err := service.RecordQueryRepositoryInterface.SelectRecordByID(ID)
	if err != nil {
		return res, err
	}

	return res, nil
}
