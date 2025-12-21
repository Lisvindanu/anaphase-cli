package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/lisvindanuu/anaphase-cli/internal/core/entity"
	"github.com/lisvindanuu/anaphase-cli/internal/core/valueobject"
)

// CustomerRepository defines the contract for customer persistence
type CustomerRepository interface {
	// Save Saves a customer entity to the repository.
	Save(ctx context.Context, customer *entity.Customer) error

	// FindByID Retrieves a customer entity by its unique ID.
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)

	// FindByEmail Retrieves a customer entity by its email address.
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error)

}
