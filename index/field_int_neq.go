package index

import (
	"github.com/juju/errgo"
)

type field_int_neq_t struct {
	field_with_int_value_t
	field_t
}

func (self *field_int_neq_t) Add(id Id, value interface{}) error {
	var err error = nil

	if id < 0 {
		return ErrorUnsaved
	}

	conn := self.idx.Conn()
	defer conn.Close()

	_, err = conn.Do("ZADD", self.DataKey(), value, id)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (self *field_int_neq_t) Del(id Id) error {
	if id < 0 {
		return ErrorUnsaved
	}

	conn := self.idx.Conn()
	defer conn.Close()

	_, err := conn.Do("ZREM", self.DataKey(), id)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (self *field_int_neq_t) Filter(op string, value interface{}, sourceKey string, targetKey string) error {
	// value is expected to be an array of 2 items
	// val := value.([]interface{})

	return errgo.New("not implemented")
}
