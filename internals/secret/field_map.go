package secret

import "errors"

type FieldType int

const (
	Text    FieldType = iota // Plain text
	Hidden                   // Secure, hidden text
	Boolean                  // Boolean (true/false)
)

// FieldValue represents a value with a specific type.
type FieldValue struct {
	Type  FieldType
	Value interface{} // Holds the actual value; expected types are string (Text, Hidden) or bool (Boolean)
}

// FieldMap is a custom type for managing fields of type map[string]FieldValue.
type FieldMap struct {
	fields map[string]FieldValue
}

// NewFieldMap initializes a new FieldMap.
func NewFieldMap() *FieldMap {
	return &FieldMap{
		fields: make(map[string]FieldValue),
	}
}

// CreateField adds a new field if it doesn't already exist.
func (fm *FieldMap) CreateField(key string, fieldValue FieldValue) error {
	if _, exists := fm.fields[key]; exists {
		return errors.New("field already exists")
	}
	fm.fields[key] = fieldValue
	return nil
}

// ReadField retrieves a field by key.
func (fm *FieldMap) ReadField(key string) (FieldValue, error) {
	value, exists := fm.fields[key]
	if !exists {
		return FieldValue{}, errors.New("field not found")
	}
	return value, nil
}

// UpdateField updates an existing field's value.
func (fm *FieldMap) UpdateField(key string, newValue FieldValue) error {
	if _, exists := fm.fields[key]; !exists {
		return errors.New("field not found")
	}
	fm.fields[key] = newValue
	return nil
}

// DeleteField removes a field by key.
func (fm *FieldMap) DeleteField(key string) error {
	if _, exists := fm.fields[key]; !exists {
		return errors.New("field not found")
	}
	delete(fm.fields, key)
	return nil
}

// ListFields returns all fields as a map.
func (fm *FieldMap) ListFields() map[string]FieldValue {
	return fm.fields
}
