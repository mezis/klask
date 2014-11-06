package index

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

type Field interface {
	Name() string
	Type() FieldType

	// Verify whether the persisted keys matchs the expected type for this field.
	// Return nil if everything is sound.
	// Check() error

	// Key to a Redis of keys containing the actual data.
	// Also a prefix for all data for this field.
	DataKey() string

	// Return a normalised value if this value is valid for this field type,
	// fails if invalid.
	// For instance passed a float64 to an int field should return the integer
	// part. Passing a string to an int field should convert the string to an
	// integer, or error if the string does not contain an integer.
	CheckValidValue(interface{}) (interface{}, error)

	// Persists the given `id` associated to the `value`.
	// The value type is dependent on the concerte implementation (given by
	// `#Type`).
	// Idempotent.
	Add(id Id, value interface{}) error

	// Removes this persisted `id`, if present. Idempotent.
	Del(id Id) error

	// Stores in the Redis `key` a set of record IDs
	// matching the operator `op` and operand `val`.
	// And error is returned if the operator, operand, or their combination have
	// unsuppoerted types or values.
	//
	// - `sourceKey`: a SET of record IDs be filtered
	// - `targetKey`: where to store the filtered set of record IDs
	Filter(op string, val interface{}, sourceKey string, targetKey string) error
}

type field_t struct {
	idx  Index
	name string
	ty   FieldType
}

func (self *field_t) Name() string {
	return self.name
}

func (self *field_t) Type() FieldType {
	return self.ty
}

func (self *field_t) keyType(key string) (string, error) {
	conn := self.idx.Conn()
	defer conn.Close()

	val, err := redis.String(conn.Do("TYPE", key))
	if err != nil {
		return "", errgo.Mask(err)
	}
	return val, nil
}

func (self *field_t) DataKey() string {
	return fmt.Sprintf("data:%s:%s", self.idx.Name(), self.Name())
}

func initField(self *field_t, idx Index, name string, ty FieldType) {
	self.idx = idx
	self.name = name
	self.ty = ty
}

func (self *field_t) init(idx Index, name string, ty FieldType) {
	self.idx = idx
	self.name = name
	self.ty = ty
}

// Factory for fields of different types and indexing methods
func newField(idx Index, name string, ty FieldType) (Field, error) {

	// TODO: validate name

	if !ty.IsValid() {
		return nil, errgo.Newf("invalid field type '%s'", ty)
	}
	switch ty {
	case FIntEq:
		field := new(field_int_eq_t)
		(&field.field_t).init(idx, name, ty)
		return field, nil
	case FIntNeq:
		field := new(field_int_neq_t)
		(&field.field_t).init(idx, name, ty)
		return field, nil
	default:
		return nil, errgo.Newf("unsupported field type '%s'", ty)
	}
}
