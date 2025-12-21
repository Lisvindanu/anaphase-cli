-- Schema for customers table

CREATE TABLE IF NOT EXISTS customers (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	-- TODO: Add domain-specific columns
);

CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at);
