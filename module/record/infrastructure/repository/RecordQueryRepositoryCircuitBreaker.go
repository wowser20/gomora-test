package repository

import (
	"github.com/afex/hystrix-go/hystrix"

	"gomora/module/record/domain/entity"
	"gomora/module/record/domain/repository"
)

// RecordQueryRepositoryCircuitBreaker holds the implementable methods for record query circuitbreaker
type RecordQueryRepositoryCircuitBreaker struct {
	repository.RecordQueryRepositoryInterface
}

// SelectRecords decorator pattern for select records repository
func (repository *RecordQueryRepositoryCircuitBreaker) SelectRecords(page *uint) ([]entity.Record, uint, error) {
	type outputData struct {
		Records    []entity.Record
		TotalCount uint
	}
	output := make(chan outputData, 1)
	errChan := make(chan error, 1)

	hystrix.ConfigureCommand("select_records", config.Settings())
	errors := hystrix.Go("select_records", func() error {
		records, totalCount, err := repository.RecordQueryRepositoryInterface.SelectRecords(page)
		if err != nil {
			errChan <- err
			return nil
		}

		result := outputData{
			Records:    records,
			TotalCount: totalCount,
		}

		output <- result
		return nil
	}, nil)

	select {
	case out := <-output:
		return out.Records, out.TotalCount, nil
	case err := <-errChan:
		return []entity.Record{}, 0, err
	case err := <-errors:
		return []entity.Record{}, 0, err
	}
}

// SelectRecordByID decorator pattern for select record repository
func (repository *RecordQueryRepositoryCircuitBreaker) SelectRecordByID(ID string) (entity.Record, error) {
	output := make(chan entity.Record, 1)
	errChan := make(chan error, 1)

	hystrix.ConfigureCommand("select_record_by_id", config.Settings())
	errors := hystrix.Go("select_record_by_id", func() error {
		record, err := repository.RecordQueryRepositoryInterface.SelectRecordByID(ID)
		if err != nil {
			errChan <- err
			return nil
		}

		output <- record
		return nil
	}, nil)

	select {
	case out := <-output:
		return out, nil
	case err := <-errChan:
		return entity.Record{}, err
	case err := <-errors:
		return entity.Record{}, err
	}
}
