package index

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

type Field interface {
	Name() string
	Type() FieldType
	Save() error

	// Redis key holding the field data
	DataKey() string

	// Destroy persisted field
	// Destroy() error
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

func (self *field_t) saveCommon() error {
	conn := self.idx.Conn()
	key := self.idx.FieldsKey()

	// TODO: this needs to be transactional
	val, err := conn.Do("HGET", key, self.name)
	if err != nil {
		return errgo.Mask(err)
	}
	if val != nil {
		if val, _ := redis.String(val, err); val != string(self.ty) {
			return errgo.Newf("field '%s' already has type '%s'", self.name, val)
		}
	}

	_, err = conn.Do("HSET", key, self.name, self.ty)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *field_t) keyType(key string) (string, error) {
	conn := self.idx.Conn()
	val, err := redis.String(conn.Do("TYPE", key))
	if err != nil {
		return "", errgo.Mask(err)
	}
	return val, nil
}

func (self *field_t) DataKey() string {
	return fmt.Sprintf("data:%s:%s", self.idx.Name(), self.Name())
}

// Factory for fields of different types and indexing methods
func newField(idx Index, name string, ty FieldType) (Field, error) {
	var err error = nil
	var result Field = nil

	// TODO: validate name

	if !ty.IsValid() {
		return nil, errgo.Newf("invalid field type '%s'", ty)
	}
	switch ty {
	case FIntEq:
		result = &field_int_eq_t{field_t{idx, name, ty}}
	case FIntNeq:
		result = &field_int_neq_t{field_t{idx, name, ty}}
	default:
		err = errgo.Newf("unsupported field type '%s'", ty)
	}
	return result, err
}
