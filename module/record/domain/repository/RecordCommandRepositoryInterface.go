package repository

import (
	"gomora/module/record/domain/entity"
	"gomora/module/record/infrastructure/repository/types"
)

// RecordCommandRepositoryInterface holds the implementable methods for record command repository
type RecordCommandRepositoryInterface interface {
	// DeleteRecord deletes a record
	DeleteRecord(ID string) error
	// InsertRecord creates a new record
	InsertRecord(data types.CreateRecord) (entity.Record, error)
	// UpdateRecord updates a record
	UpdateRecord(data types.UpdateRecord) error
}
