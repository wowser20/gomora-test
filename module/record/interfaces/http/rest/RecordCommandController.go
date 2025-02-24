package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"gomora/interfaces/http/rest/viewmodels"
	"gomora/internal/errors"
	apiError "gomora/internal/errors"
	"gomora/module/record/application"
	serviceTypes "gomora/module/record/infrastructure/service/types"
	types "gomora/module/record/interfaces/http"
)

// RecordCommandController request controller for record command
type RecordCommandController struct {
	application.RecordCommandServiceInterface
}

// CreateRecord request handler to create record
func (controller *RecordCommandController) CreateRecord(w http.ResponseWriter, r *http.Request) {
	var request types.CreateRecordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := viewmodels.HTTPResponseVM{
			Status:    http.StatusBadRequest,
			Success:   false,
			Message:   "Invalid payload request.",
			ErrorCode: apiError.InvalidRequestPayload,
		}

		response.JSON(w)
		return
	}

	// validate request
	err := types.Validate.Struct(request)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		if len(errors) > 0 {
			response := viewmodels.HTTPResponseVM{
				Status:    http.StatusBadRequest,
				Success:   false,
				Message:   types.ValidationErrors[errors[0].StructNamespace()],
				ErrorCode: apiError.InvalidPayload,
			}

			response.JSON(w)
			return
		}

		response := viewmodels.HTTPResponseVM{
			Status:    http.StatusBadRequest,
			Success:   false,
			Message:   "Invalid payload request.",
			ErrorCode: apiError.InvalidRequestPayload,
		}

		response.JSON(w)
		return
	}

	record := serviceTypes.CreateRecord{
		ID:   request.ID,
		Data: request.Data,
	}

	res, err := controller.RecordCommandServiceInterface.CreateRecord(context.TODO(), record)
	if err != nil {
		var httpCode int
		var errorMsg string

		switch err.Error() {
		case errors.DatabaseError:
			httpCode = http.StatusInternalServerError
			errorMsg = "Error occurred while saving record."
		case errors.DuplicateRecord:
			httpCode = http.StatusConflict
			errorMsg = "Record ID already exist."
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
		Message: "Successfully created record.",
		Data: &types.CreateRecordResponse{
			ID:        res.ID,
			Data:      res.Data,
			CreatedAt: time.Now().Unix(),
		},
	}

	response.JSON(w)
}

// DeleteRecord deletes a record
func (controller *RecordCommandController) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	if len(ID) == 0 {
		response := viewmodels.HTTPResponseVM{
			Status:    http.StatusBadRequest,
			Success:   false,
			Message:   "ID is required.",
			ErrorCode: apiError.InvalidRequestPayload,
		}

		response.JSON(w)
		return
	}

	err := controller.RecordCommandServiceInterface.DeleteRecord(context.TODO(), ID)
	if err != nil {
		var httpCode int
		var errorMsg string

		switch err.Error() {
		case errors.DatabaseError:
			httpCode = http.StatusInternalServerError
			errorMsg = "Error occurred while deleting record."
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
		Message: "Successfully deleted record.",
	}

	response.JSON(w)
}

// GenerateToken request handler to generate token
func (controller *RecordCommandController) GenerateToken(w http.ResponseWriter, r *http.Request) {
	token, err := controller.RecordCommandServiceInterface.GenerateToken(context.TODO())
	if err != nil {
		var httpCode int
		var errorMsg string

		switch err.Error() {
		case errors.DatabaseError:
			httpCode = http.StatusInternalServerError
			errorMsg = "Error occurred while generating token."
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
		Message: "Successfully generated token.",
		Data: &types.GenerateTokenResponse{
			AccessToken: token,
		},
	}

	response.JSON(w)
}

// UpdateRecord updates a record
func (controller *RecordCommandController) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	if len(ID) == 0 {
		response := viewmodels.HTTPResponseVM{
			Status:    http.StatusBadRequest,
			Success:   false,
			Message:   "ID is required.",
			ErrorCode: apiError.InvalidRequestPayload,
		}

		response.JSON(w)
		return
	}

	var request types.UpdateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := viewmodels.HTTPResponseVM{
			Status:    http.StatusBadRequest,
			Success:   false,
			Message:   "Invalid payload request.",
			ErrorCode: apiError.InvalidRequestPayload,
		}

		response.JSON(w)
		return
	}

	err := controller.RecordCommandServiceInterface.UpdateRecord(context.TODO(), serviceTypes.UpdateRecord{
		ID:   ID,
		Data: request.Data,
	})
	if err != nil {
		var httpCode int
		var errorMsg string

		switch err.Error() {
		case errors.DatabaseError:
			httpCode = http.StatusInternalServerError
			errorMsg = "Error occurred while updating record."
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
		Message: "Successfully updated record.",
	}

	response.JSON(w)
}
