package application

import (
	"context"

	"gomora/module/record/domain/entity"
)

// RecordQueryServiceInterface holds the implementable methods for the record query service
type RecordQueryServiceInterface interface {
	// GetRecords gets all records
	GetRecords(ctx context.Context, page *uint) ([]entity.Record, uint, error)
	// GetRecordByID gets a record by its ID
	GetRecordByID(ctx context.Context, ID string) (entity.Record, error)
}
