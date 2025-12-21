package http

// CreateCustomerRequest represents HTTP request to create customer
type CreateCustomerRequest struct {
	// TODO: Add fields based on domain entity
}

// UpdateCustomerRequest represents HTTP request to update customer
type UpdateCustomerRequest struct {
	// TODO: Add fields based on domain entity
}

// CustomerResponse represents HTTP response with customer data
type CustomerResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	// TODO: Add fields based on domain entity
}

// ErrorResponse represents standard error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}
