package valueobject

// Email is a value object
type Email struct {
	Value string // The email address
}

// NewEmail creates a new Email
func NewEmail(value string) (Email, error) {
	v := Email{
		Value: value,
	}

	if err := v.Validate(); err != nil {
		return Email{}, err
	}

	return v, nil
}

// Validate validates the Email
func (v Email) Validate() error {
	// Value must be a valid email format and not empty
	// TODO: Add validation logic
	return nil
}
