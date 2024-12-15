package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/neox5/openk/internal/kms"
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
type derivationV1Handler struct {
	storage storage.MiniStorageBackend
}

func newDerivationV1Handler(storage storage.MiniStorageBackend) *derivationV1Handler {
	if storage == nil {
		panic("storage cannot be nil")
	}
	return &derivationV1Handler{storage: storage}
}

func (h *derivationV1Handler) storeParams(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req storeParamsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, r, ValidationError("body", "invalid JSON"))
		return
	}

	// Validate request
	if err := validateDerivationRequest(&req); err != nil {
		WriteError(w, r, err)
		return
	}

	// Create derivation parameters
	params, err := kms.NewKeyDerivation(req.Username, req.Iterations)
	if err != nil {
		switch err {
		case kms.ErrUsernameEmpty, kms.ErrUsernameLength, kms.ErrUsernameInvalid:
			WriteError(w, r, ValidationError("username", err.Error()))
			return
		case kms.ErrIterationsInvalid:
			WriteError(w, r, ValidationError("iterations", err.Error()))
			return
		default:
			WriteError(w, r, InternalError().WithDetail(err.Error()))
			return
		}
	}

	// Store parameters
	stored, err := h.storage.StoreDerivationParams(r.Context(), params)
	if err != nil {
		switch err {
		case storage.ErrParamsNil:
			WriteError(w, r, ValidationError("params", "cannot be nil"))
			return
		case storage.ErrParamsNotFound:
			WriteError(w, r, NotFoundError("derivation parameters", params.Username))
			return
		default:
			WriteError(w, r, InternalError().WithDetail(err.Error()))
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

	WriteJSON(w, http.StatusCreated, resp)
}

func (h *derivationV1Handler) getParams(w http.ResponseWriter, r *http.Request) {
	// Get username from path parameter using Go 1.22's PathValue
	username := r.PathValue("username")
	if username == "" {
		WriteError(w, r, ValidationError("username", "cannot be empty"))
		return
	}

	// Get parameters from storage
	params, err := h.storage.GetDerivationParams(r.Context(), username)
	if err != nil {
		switch err {
		case kms.ErrUsernameEmpty:
			WriteError(w, r, ValidationError("username", "cannot be empty"))
			return
		case storage.ErrParamsNotFound:
			WriteError(w, r, NotFoundError("derivation parameters", username))
			return
		default:
			WriteError(w, r, InternalError().WithDetail(err.Error()))
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

	WriteJSON(w, http.StatusOK, resp)
}

func validateDerivationRequest(req *storeParamsRequest) error {
	if req.Username == "" {
		return ValidationError("username", "cannot be empty")
	}
	if len(req.Username) > kms.MaxUsernameLen {
		return ValidationError("username", "exceeds maximum length")
	}
	if req.Iterations < kms.MinIterations {
		return ValidationError("iterations", "below minimum value")
	}
	return nil
}
