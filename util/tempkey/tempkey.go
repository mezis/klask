package tempkey

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

const cTempKeyPrefix = "tempkey"

type Keys interface {
	// Return the name of an empty key
	Get() (string, error)

	// Release a key obtained by Get()
	// It is not an error if the key was never obtained, or was released
	// previously.
	// Connection errors are still reported.
	Release(string) error

	// Release all keys created with Get()
	// It is an error if any removal fails, although each key will be attempted.
	Clear() error
}

type ConnFactory func() redis.Conn

type keys_t struct {
	factory ConnFactory
	keys    map[string]bool
}

func New(fact ConnFactory) Keys {
	rv := new(keys_t)
	rv.factory = fact
	rv.keys = make(map[string]bool)
	return rv
}

func (self *keys_t) Get() (string, error) {
	conn := self.factory()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("INCR", cTempKeyPrefix))
	if err != nil {
		return "", errgo.Mask(err)
	}

	key := fmt.Sprintf("%s:%d", cTempKeyPrefix, val)
	self.keys[key] = true
	return key, nil
}

func (self *keys_t) Release(key string) error {
	if _, ok := self.keys[key]; !ok {
		return nil
	}

	conn := self.factory()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return errgo.Mask(err)
	}
	delete(self.keys, key)

	return nil
}

func (self *keys_t) Clear() error {
	var err error = nil
	for key, ok := range self.keys {
		if !ok {
			continue
		}
		e := self.Release(key)
		if e != nil {
			err = errgo.Mask(e)
		}
	}
	return err
}

// func NewSlice(conn redis.Conn, size int) (Slice, error) {
// 	slice := make([]Key, size)

// 	for k, _ := range slice {
// 		key, err := New(conn)
// 		if err != nil {
// 			return nil, errgo.Mask(err)
// 		}
// 		slice[k] = key
// 	}

// 	return slice, nil
// }

// func (self *key_t) Name() string {
// 	return self.key
// }

// func (self *key_t) Clear() error {
// 	_, err := self.conn.Do("DEL", self.key)
// 	if err != nil {
// 		return errgo.Mask(err)
// 	}
// 	return nil
// }

// func (self Slice) Clear() error {
// 	var err error = nil
// 	for _, key := range self {
// 		if e := key.Clear(); e != nil {
// 			err = errgo.Mask(e)
// 		}
// 	}
// 	return err
// }
