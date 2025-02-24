package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"gomora/infrastructures/database/mysql/types"
	apiError "gomora/internal/errors"
	"gomora/module/record/domain/entity"
)

// RecordQueryRepository handles the record query repository logic
type RecordQueryRepository struct {
	types.MySQLDBHandlerInterface
}

// SelectRecords select all records
func (repository *RecordQueryRepository) SelectRecords(page *uint) ([]entity.Record, uint, error) {
	var record entity.Record
	var records []entity.Record

	stmt := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC", record.GetModelName())

	// get total count
	var counter struct {
		Total uint `json:"total"`
	}

	totalCountStmt := strings.ReplaceAll(stmt, "SELECT *", "SELECT COUNT(*) as total")
	err := repository.QueryRow(totalCountStmt, map[string]interface{}{}, &counter)
	if err != nil {
		log.Println(err)
		return nil, 0, errors.New(apiError.DatabaseError)
	}

	// if page is set
	if page != nil {
		if *page > 0 {
			var limit uint = 10
			offset := limit * (*page - 1)

			stmt = fmt.Sprintf("%s LIMIT %d OFFSET %d", stmt, limit, offset)
		}
	}

	err = repository.Query(stmt, map[string]interface{}{}, &records)
	if err != nil {
		log.Println(err)
		return nil, 0, errors.New(apiError.DatabaseError)
	} else if len(records) == 0 {
		return nil, 0, errors.New(apiError.MissingRecord)
	}

	return records, counter.Total, nil
}

// SelectRecordByID select a record by id
func (repository *RecordQueryRepository) SelectRecordByID(ID string) (entity.Record, error) {
	var record entity.Record

	stmt := fmt.Sprintf("SELECT * FROM %s WHERE id=:id", record.GetModelName())
	err := repository.QueryRow(stmt, map[string]interface{}{
		"id": ID,
	}, &record)
	if err != nil {
		if err == sql.ErrNoRows {
			return record, errors.New(apiError.MissingRecord)
		}

		return record, errors.New(apiError.DatabaseError)
	}

	return record, nil
}
