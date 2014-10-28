package index

import (
	"errors"
)

type Field interface {
	Name() string
	Type() FieldType
	Save(*Index) error
	Init(string, FieldType) error
}

type field_t struct {
	name string
	ty   FieldType
}

type field_int_eq_t struct {
	field_t
}

type field_int_neq_t struct {
	field_t
}

func (self *field_t) Init(name string, ty FieldType) error {
	return nil
}

func (self *field_t) Name() string {
	return self.name
}

func (self *field_t) Type() FieldType {
	return self.ty
}

func (self *field_int_eq_t) Save(idx *Index) error {
	return errors.New("not implemented")
}

func (self *field_int_neq_t) Save(idx *Index) error {
	return errors.New("not implemented")
}

// Factory for fields of different types and indexing methods
func NewField(name string, ty FieldType) (Field, error) {
	var err error
	var result Field

	if !ty.IsValid() {
		return nil, errors.New("invalid field type")
	}

	switch ty {
	case FIntEq:
		result = new(field_int_eq_t)
	case FIntNeq:
		result = new(field_int_neq_t)
	default:
		err = errors.New("unsupported field type")
		return nil, err
	}
	err = result.Init(name, ty)
	return result, err
}
