package application

import (
	"context"

	"gomora/module/record/domain/entity"
	"gomora/module/record/infrastructure/service/types"
)

// RecordCommandServiceInterface holds the implementable methods for the record command service
type RecordCommandServiceInterface interface {
	// CreateRecord creates a new record
	CreateRecord(ctx context.Context, data types.CreateRecord) (entity.Record, error)
	// DeleteRecord deletes a record
	DeleteRecord(ctx context.Context, ID string) error
	// GenerateToken generates a jwt token
	GenerateToken(ctx context.Context) (string, error)
	// UpdateRecord updates a record
	UpdateRecord(ctx context.Context, data types.UpdateRecord) error
}
