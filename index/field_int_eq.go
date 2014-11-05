package index

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

type field_int_eq_t struct {
	field_with_int_value_t
	field_t
}

func (self *field_int_eq_t) valueKey(val int64) string {
	return fmt.Sprintf("%s:%d", self.DataKey(), val)
}

func (self *field_int_eq_t) Add(id Id, value interface{}) error {
	// TODO: check type properly
	var (
		err error = nil
		val int64 = value.(int64)
	)

	val = value.(int64)
	if id < 0 {
		return ErrorUnsaved
	}

	conn := self.idx.Conn()
	defer conn.Close()

	// TODO: wrap these in a Lua script
	_, err = conn.Do("SADD", self.valueKey(val), id)
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = conn.Do("HSET", self.DataKey(), id, val)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (self *field_int_eq_t) Del(id Id) error {
	var (
		err error = nil
		val int64
	)

	if id < 0 {
		return ErrorUnsaved
	}

	conn := self.idx.Conn()
	defer conn.Close()

	// TODO: wrap these in a Lua script
	val, err = redis.Int64(conn.Do("HGET", self.DataKey(), id))
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = conn.Do("SREM", self.valueKey(val), id)
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = conn.Do("HDEL", self.DataKey(), id)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (self *field_int_eq_t) eqFilter(valueKey string, sourceKey string, targetKey string) error {
	conn := self.idx.Conn()
	defer conn.Close()

	_, err := conn.Do("SINTERSTORE", targetKey, sourceKey, valueKey)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *field_int_eq_t) neqFilter(valueKey string, sourceKey string, targetKey string) error {
	conn := self.idx.Conn()
	defer conn.Close()

	_, err := conn.Do("SDIFFSTORE", targetKey, sourceKey, valueKey)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *field_int_eq_t) Filter(op string, value interface{}, sourceKey string, targetKey string) error {
	// TODO: better type handling
	var (
		val int64 = value.(int64)
		err error = nil
	)

	conn := self.idx.Conn()
	defer conn.Close()

	valueKey := self.valueKey(val)
	switch op {
	case "eq":
		err = self.eqFilter(valueKey, sourceKey, targetKey)
	case "neq":
		err = self.neqFilter(valueKey, sourceKey, targetKey)
	default:
		return errgo.Newf("bad operator '%s'", op)
	}
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
