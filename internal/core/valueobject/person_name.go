package valueobject

// PersonName is a value object
type PersonName struct {
	Value string // The full name of a person
}

// NewPersonName creates a new PersonName
func NewPersonName(value string) (PersonName, error) {
	v := PersonName{
		Value: value,
	}

	if err := v.Validate(); err != nil {
		return PersonName{}, err
	}

	return v, nil
}

// Validate validates the PersonName
func (v PersonName) Validate() error {
	// Value must not be empty and should contain only letters and spaces
	// TODO: Add validation logic
	return nil
}
