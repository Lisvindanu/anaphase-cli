package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lisvindanu/anaphase-cli/internal/core/entity"
	"github.com/lisvindanu/anaphase-cli/internal/core/port"
	"github.com/lisvindanu/anaphase-cli/internal/core/valueobject"
)

type customerRepository struct {
	db *pgxpool.Pool
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *pgxpool.Pool) port.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

// Save saves a customer to the repository
func (r *customerRepository) Save(ctx context.Context, entity *entity.Customer) error {
	query := `
		INSERT INTO customers (id, created_at, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET updated_at = $3
	`

	_, err := r.db.Exec(ctx, query,
		entity.ID,
		entity.CreatedAt,
		entity.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("save customer: %w", err)
	}

	return nil
}

// FindByID retrieves a customer by ID
func (r *customerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	var result entity.Customer

	query := `
		SELECT id, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("find customer: %w", err)
	}

	return &result, nil
}

// FindByEmail retrieves a customer by email address
func (r *customerRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error) {
	// TODO: Implement FindByEmail
	return nil, fmt.Errorf("not implemented")
}

