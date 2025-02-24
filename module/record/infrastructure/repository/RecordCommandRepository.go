package repository

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"gomora/infrastructures/database/mysql/types"
	apiError "gomora/internal/errors"
	"gomora/module/record/domain/entity"
	repositoryTypes "gomora/module/record/infrastructure/repository/types"
)

// RecordCommandRepository handles the record command repository logic
type RecordCommandRepository struct {
	types.MySQLDBHandlerInterface
}

// DeleteRecord deletes a record
func (repository *RecordCommandRepository) DeleteRecord(ID string) error {
	record := entity.Record{
		ID: ID,
	}

	stmt := fmt.Sprintf("DELETE FROM %s WHERE id=:id", record.GetModelName())
	res, err := repository.MySQLDBHandlerInterface.Execute(stmt, record)
	if err != nil {
		return errors.New(apiError.DatabaseError)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New(apiError.DatabaseError)
	}

	return nil
}

// InsertRecord creates a new record
func (repository *RecordCommandRepository) InsertRecord(data repositoryTypes.CreateRecord) (entity.Record, error) {
	record := entity.Record{
		ID:   data.ID,
		Data: data.Data,
	}

	stmt := fmt.Sprintf("INSERT INTO %s (id, data) VALUES (:id, :data)", record.GetModelName())
	_, err := repository.MySQLDBHandlerInterface.Execute(stmt, record)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return entity.Record{}, errors.New(apiError.DuplicateRecord)
		}
		return entity.Record{}, errors.New(apiError.DatabaseError)
	}

	return record, nil
}

// UpdateRecord updates a record
func (repository *RecordCommandRepository) UpdateRecord(data repositoryTypes.UpdateRecord) error {
	record := entity.Record{
		ID:   data.ID,
		Data: data.Data,
	}

	stmt := fmt.Sprintf("UPDATE %s SET data=:data WHERE id=:id", record.GetModelName())
	_, err := repository.MySQLDBHandlerInterface.Execute(stmt, record)
	if err != nil {
		return errors.New(apiError.DatabaseError)
	}

	return nil
}
