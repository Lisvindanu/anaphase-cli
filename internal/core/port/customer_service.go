package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/lisvindanu/anaphase-cli/internal/core/entity"
	"github.com/lisvindanu/anaphase-cli/internal/core/valueobject"
)

// CustomerService defines the contract for customer business logic
type CustomerService interface {
	// RegisterCustomer Registers a new customer with the provided name and email.
	RegisterCustomer(ctx context.Context, name valueobject.PersonName, email valueobject.Email) (*entity.Customer, error)

	// UpdateCustomerDetails Updates the name and email for an existing customer.
	UpdateCustomerDetails(ctx context.Context, customerID uuid.UUID, name valueobject.PersonName, email valueobject.Email) (*entity.Customer, error)

	// PlaceOrder Allows a customer to place an order for specified products and quantities. Returns the new order's ID.
	PlaceOrder(ctx context.Context, customerID uuid.UUID, productIDs []uuid.UUID, quantities []int) (uuid.UUID, error)
}
