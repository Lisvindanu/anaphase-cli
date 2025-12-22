package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lisvindanu/anaphase-cli/internal/core/valueobject"
)

// Common errors
var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrInvalidCustomer  = errors.New("invalid customer")
)

// Customer is an aggregate root
type Customer struct {
	ID        uuid.UUID              // Unique identifier for the customer
	Name      valueobject.PersonName // Customer's full name
	Email     valueobject.Email      // Customer's email address
	CreatedAt time.Time              // Timestamp when the customer record was created
	UpdatedAt time.Time              // Timestamp when the customer record was last updated
}

// NewCustomer creates a new customer
func NewCustomer() *Customer {
	return &Customer{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate validates the customer
func (e *Customer) Validate() error {
	if e.ID == uuid.Nil {
		return ErrInvalidCustomer
	}
	// Validate ID: cannot be nil
	// Validate Name: must be a valid PersonName value object
	// Validate Email: must be a valid Email value object
	// Validate CreatedAt: cannot be zero
	// Validate UpdatedAt: cannot be zero
	return nil
}

// UpdateName Updates the customer's name.
func (c *Customer) UpdateName(newName valueobject.PersonName) error {
	// TODO: Implement business logic
	return nil
}

// UpdateEmail Updates the customer's email address.
func (c *Customer) UpdateEmail(newEmail valueobject.Email) error {
	// TODO: Implement business logic
	return nil
}
