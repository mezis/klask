package index

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
)

type Index interface {
	json.Marshaler
	json.Unmarshaler

	// Return the index name
	Name() string

	// Checks the database for index existence
	Exists() (bool, error)

	// Populates the structure from persisted databse records.
	// Errors if the persisted version cannot be read, is inconsistent, or does not
	// exist.
	Load() error
	Save() error
	Destroy() error
	Fields() Fieldset
}

type index_t struct {
	pool   *redis.Pool // the pool of Redis connections to use
	name   string      // the index name, used as a prefix in Redis
	dirty  bool        // whether changes were made that need to be persisted
	fields Fieldset    // maps field names to how they should be indexed
}

func New(name string, pool *redis.Pool) (Index, error) {
	self := new(index_t)
	// TODO: validate name
	self.pool = pool
	self.name = name
	self.fields = make(Fieldset)
	return self, nil
}

func (self *index_t) Name() string {
	return self.name
}

func (self *index_t) Exists() (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()

	res, err := redis.Bool(conn.Do("SISMEMBER", "indices", self.name))
	if err != nil {
		return false, err
	}

	return res, nil
}

func (self *index_t) Load() error {
	return errors.New("not implemented")
}

// TODO: make this transactional, using version numbers/UUID and
// a LUA script
func (self *index_t) Save() error {
	conn := self.pool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("SADD", "indices", self.name))
	if err != nil {
		return err
	}

	return nil
}

func (self *index_t) Destroy() error {
	return errors.New("not implemented")
}

func (self *index_t) Fields() Fieldset {
	return self.fields
}
