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

	// Release one or more keys obtained by Get().
	// All arguments should be strings.
	// It is not an error if a key was never obtained, or was released
	// previously.
	// Connection and argument errors are still reported.
	Release(...interface{}) error

	// Release all keys created with Get()
	// It is an error if any removal fails, although each key will be attempted.
	Clear() error
}

type ConnFactory func() redis.Conn

type keys_t struct {
	factory ConnFactory
	keys    map[string]bool
}

//

func New(fact ConnFactory) Keys {
	rv := new(keys_t)
	rv.factory = fact
	rv.keys = make(map[string]bool)
	return rv
}

//

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

//

func (self *keys_t) Release(keys ...interface{}) error {
	to_delete := make([]interface{}, 0, len(keys))

	for _, key := range keys {
		switch k := key.(type) {
		case string:
			if _, ok := self.keys[k]; !ok {
				continue
			}
			to_delete = append(to_delete, k)
		default:
			return errgo.New("arguments should be strings")
		}
	}

	conn := self.factory()
	defer conn.Close()

	_, err := conn.Do("DEL", to_delete...)
	if err != nil {
		return errgo.Mask(err)
	}
	for _, key := range to_delete {
		delete(self.keys, key.(string))
	}

	return nil
}

//

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
