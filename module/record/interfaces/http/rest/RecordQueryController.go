package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"gomora/interfaces/http/rest/viewmodels"
	"gomora/internal/errors"
	"gomora/module/record/application"
	types "gomora/module/record/interfaces/http"
)

// RecordQueryController request controller for record query
type RecordQueryController struct {
	application.RecordQueryServiceInterface
}

// GetRecords retrieves all records
func (controller *RecordQueryController) GetRecords(w http.ResponseWriter, r *http.Request) {
	// optional
	var page *uint
	pageStr := r.URL.Query().Get("page")
	if len(pageStr) > 0 {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt == 0 {
			page = nil
		}
		pageUint := uint(pageInt)
		page = &pageUint
	}

	res, total, err := controller.RecordQueryServiceInterface.GetRecords(context.TODO(), page)
	if err != nil {
		var httpCode int
		var errorMsg string

		switch err.Error() {
		default:
			httpCode = http.StatusInternalServerError
			errorMsg = "Database error."
		}

		response := viewmodels.HTTPResponseVM{
			Status:    httpCode,
			Success:   false,
			Message:   errorMsg,
			ErrorCode: err.Error(),
		}

		response.JSON(w)
		return
	}

	records := []types.GetRecordResponse{}

	for _, record := range res {
		records = append(records, types.GetRecordResponse{
			ID:        record.ID,
			Data:      record.Data,
			CreatedAt: record.CreatedAt.Unix(),
		})
	}

	response := viewmodels.HTTPResponseVM{
		Status:  http.StatusOK,
		Success: true,
		Message: "Successfully fetched records.",
		Data: &types.GetPaginatedRecordsResponse{
			Records: records,
			Total:   total,
		},
	}

	response.JSON(w)
}

// GetRecordByID retrieves the tenant id from the rest request
func (controller *RecordQueryController) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	recordID := chi.URLParam(r, "id")

	if len(recordID) == 0 {
		response := viewmodels.HTTPResponseVM{
			Status:    http.StatusBadRequest,
			Success:   false,
			Message:   "Invalid record ID",
			ErrorCode: errors.InvalidRequestPayload,
		}

		response.JSON(w)
		return
	}

	res, err := controller.RecordQueryServiceInterface.GetRecordByID(context.TODO(), recordID)
	if err != nil {
		var httpCode int
		var errorMsg string

		switch err.Error() {
		case errors.DatabaseError:
			httpCode = http.StatusInternalServerError
			errorMsg = "Error while fetching record."
		case errors.MissingRecord:
			httpCode = http.StatusNotFound
			errorMsg = "No record found."
		default:
			httpCode = http.StatusInternalServerError
			errorMsg = "Please contact technical support."
		}

		response := viewmodels.HTTPResponseVM{
			Status:    httpCode,
			Success:   false,
			Message:   errorMsg,
			ErrorCode: err.Error(),
		}

		response.JSON(w)
		return
	}

	response := viewmodels.HTTPResponseVM{
		Status:  http.StatusOK,
		Success: true,
		Message: "Record successfully fetched.",
		Data: &types.GetRecordResponse{
			ID:        res.ID,
			Data:      res.Data,
			CreatedAt: res.CreatedAt.Unix(),
		},
	}

	response.JSON(w)
}
