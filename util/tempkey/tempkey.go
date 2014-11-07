package tempkey

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

const cTempKeyPrefix = "tempkey"

type Key interface {
	Name() string
	Clear() error
}

type Slice []Key

type key_t struct {
	conn redis.Conn
	key  string
}

func New(conn redis.Conn) (Key, error) {
	val, err := redis.Int64(conn.Do("INCR", cTempKeyPrefix))
	if err != nil {
		return nil, errgo.Mask(err)
	}

	key := fmt.Sprintf("%s:%d", cTempKeyPrefix, val)
	return &key_t{conn, key}, nil
}

func NewSlice(conn redis.Conn, size int) (Slice, error) {
	slice := make([]Key, size)

	for k, _ := range slice {
		key, err := New(conn)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		slice[k] = key
	}

	return slice, nil
}

func (self *key_t) Name() string {
	return self.key
}

func (self *key_t) Clear() error {
	_, err := self.conn.Do("DEL", self.key)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self Slice) Clear() error {
	var err error = nil
	for _, key := range self {
		if e := key.Clear(); e != nil {
			err = errgo.Mask(e)
		}
	}
	return err
}

func (self Slice) Names() []string {
	names := make([]string, len(self))
	for k, key := range self {
		names[k] = key.Name()
	}
	return names
}
