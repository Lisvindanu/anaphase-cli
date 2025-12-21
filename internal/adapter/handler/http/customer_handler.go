package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/lisvindanuu/anaphase-cli/internal/core/port"
)

// CustomerHandler handles HTTP requests for customer domain
type CustomerHandler struct {
	service port.CustomerService
	logger  *slog.Logger
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(service port.CustomerService, logger *slog.Logger) *CustomerHandler {
	return &CustomerHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers all routes for this handler
func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
	r.Route("/customers", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// Create creates a new customer
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// TODO: Call service to create entity
	_ = ctx
	_ = req

	h.respondJSON(w, http.StatusCreated, CustomerResponse{})
}

// GetByID retrieves customer by ID
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	// TODO: Call service to get entity
	_ = ctx
	_ = uuid

	h.respondJSON(w, http.StatusOK, CustomerResponse{})
}

// Update updates an existing customer
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	var req UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// TODO: Call service to update entity
	_ = ctx
	_ = uuid
	_ = req

	h.respondJSON(w, http.StatusOK, CustomerResponse{})
}

// Delete removes a customer
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	// TODO: Call service to delete entity
	_ = ctx
	_ = uuid

	w.WriteHeader(http.StatusNoContent)
}

// respondJSON sends a JSON response
func (h *CustomerHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// respondError sends an error response
func (h *CustomerHandler) respondError(w http.ResponseWriter, status int, message string, err error) {
	h.logger.Error(message, "error", err)
	h.respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
