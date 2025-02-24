package repository

import (
	"gomora/module/record/domain/entity"
)

// RecordQueryRepositoryInterface holds the implementable method for record query repository
type RecordQueryRepositoryInterface interface {
	// SelectRecords gets all records
	SelectRecords(page *uint) ([]entity.Record, uint, error)
	// SelectRecordByID gets a record by its ID
	SelectRecordByID(ID string) (entity.Record, error)
}
