package index

import (
	"github.com/juju/errgo"
)

type field_zset_store_t struct {
	field_t
}

func (self *field_zset_store_t) Add(id Id, value interface{}) error {
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

func (self *field_zset_store_t) Del(id Id) error {
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
