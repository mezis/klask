package index

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

const tempKeyPrefix = "tempkey"

type Tempkey interface {
	Get() string
	Clear() error
}

type key_t struct {
	idx Index
	key string
}

func (self *index_t) NewTempKey() (Tempkey, error) {
	conn := self.Conn()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("INCR", tempKeyPrefix))
	if err != nil {
		return nil, errgo.Mask(err)
	}

	key := fmt.Sprintf("%s:%d", tempKeyPrefix, val)
	return &key_t{self, key}, nil
}

func (self *key_t) Get() string {
	return self.key
}

func (self *key_t) Clear() error {
	conn := self.idx.Conn()
	defer conn.Close()

	_, err := conn.Do("DEL", self.key)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
