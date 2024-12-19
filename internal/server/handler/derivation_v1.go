package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/neox5/openk/internal/kms"
	"github.com/neox5/openk/internal/server/httperror"
	"github.com/neox5/openk/internal/storage"
)

// DTOs
type storeParamsRequest struct {
	Username   string `json:"username" validate:"required"`
	Iterations int    `json:"iterations" validate:"required,min=100000"`
}

type paramsResponse struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Iterations int       `json:"iterations"`
	CreatedAt  time.Time `json:"created_at"`
}

// Handler
type DerivationV1Handler struct {
	storage storage.MiniStorageBackend
}

func NewDerivationV1Handler(storage storage.MiniStorageBackend) *DerivationV1Handler {
	if storage == nil {
		panic("storage cannot be nil")
	}
	return &DerivationV1Handler{storage: storage}
}

func (h *DerivationV1Handler) StoreParams(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req storeParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.WriteError(w, r, httperror.ValidationError("body", "invalid JSON"))
		return
	}

	// Validate request
	if err := validateDerivationRequest(&req); err != nil {
		httperror.WriteError(w, r, err)
		return
	}

	// Create derivation parameters
	params, err := kms.NewKeyDerivation(req.Username, req.Iterations)
	if err != nil {
		switch err {
		case kms.ErrUsernameEmpty, kms.ErrUsernameLength, kms.ErrUsernameInvalid:
			httperror.WriteError(w, r, httperror.ValidationError("username", err.Error()))
			return
		case kms.ErrIterationsInvalid:
			httperror.WriteError(w, r, httperror.ValidationError("iterations", err.Error()))
			return
		default:
			httperror.WriteError(w, r, httperror.InternalError())
			return
		}
	}

	// Store parameters
	stored, err := h.storage.StoreDerivationParams(r.Context(), params)
	if err != nil {
		switch err {
		case storage.ErrParamsNil:
			httperror.WriteError(w, r, httperror.ValidationError("params", "cannot be nil"))
			return
		case storage.ErrParamsNotFound:
			httperror.WriteError(w, r, httperror.NotFoundError("derivation parameters", params.Username))
			return
		default:
			httperror.WriteError(w, r, httperror.InternalError())
			return
		}
	}

	// Return success response
	resp := &paramsResponse{
		ID:         stored.ID,
		Username:   stored.Username,
		Iterations: stored.Iterations,
		CreatedAt:  stored.CreatedAt,
	}

	httperror.WriteJSON(w, http.StatusCreated, resp)
}

func (h *DerivationV1Handler) GetParams(w http.ResponseWriter, r *http.Request) {
	// Get username from path parameter
	username := r.PathValue("username")
	if username == "" {
		httperror.WriteError(w, r, httperror.ValidationError("username", "cannot be empty"))
		return
	}

	// Get parameters from storage
	params, err := h.storage.GetDerivationParams(r.Context(), username)
	if err != nil {
		switch err {
		case kms.ErrUsernameEmpty:
			httperror.WriteError(w, r, httperror.ValidationError("username", "cannot be empty"))
			return
		case storage.ErrParamsNotFound:
			httperror.WriteError(w, r, httperror.NotFoundError("derivation parameters", username))
			return
		default:
			httperror.WriteError(w, r, httperror.InternalError())
			return
		}
	}

	// Convert to response DTO
	resp := &paramsResponse{
		ID:         params.ID,
		Username:   params.Username,
		Iterations: params.Iterations,
		CreatedAt:  params.CreatedAt,
	}

	httperror.WriteJSON(w, http.StatusOK, resp)
}

func validateDerivationRequest(req *storeParamsRequest) error {
	if req.Username == "" {
		return httperror.ValidationError("username", "cannot be empty")
	}
	if len(req.Username) > kms.MaxUsernameLen {
		return httperror.ValidationError("username", "exceeds maximum length")
	}
	if req.Iterations < kms.MinIterations {
		return httperror.ValidationError("iterations", "below minimum value")
	}
	return nil
}
