package index

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

type Index interface {
	json.Marshaler
	json.Unmarshaler

	// Return the index name
	Name() string

	// A fresh connection from the underlying pool
	// to be closed after usage
	Conn() redis.Conn

	// Checks the database for index existence
	Exists() (bool, error)

	// Populates the structure from persisted databse records.
	// Errors if the persisted version cannot be read, is inconsistent, or does not
	// exist.
	Load() error

	// Also save fields
	Save() error
	Destroy() error
	Fields() Fieldset

	// The Redis key where the list of fields is stored
	FieldsKey() string
}

// Allocate and initialize an Index
func New(name string, pool *redis.Pool) (Index, error) {
	self := new(index_t)
	// TODO: validate name
	self.pool = pool
	self.name = name
	self.fields = make(Fieldset)
	return self, nil
}

// Wrap in an Enumerable interface?
// index.Iter(pool).Each()

func Each(pool *redis.Pool) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		conn := pool.Get()
		defer conn.Close()

		names, err := redis.Strings(conn.Do("SMEMBERS", cIndicesKey))
		if err != nil {
			ch <- errgo.Mask(err)
			return
		}
		for _, name := range names {
			fmt.Println(name)
			idx, err := New(name, pool)
			if err != nil {
				ch <- errgo.Mask(err)
				return
			}
			err = idx.Load()
			if err != nil {
				ch <- errgo.Mask(err)
				return
			}
			ch <- idx
		}
		close(ch)
	}()
	return ch
}

////////////////////////////////////////////////////////////////////////////////

type index_t struct {
	pool   *redis.Pool // the pool of Redis connections to use
	name   string      // the index name, used as a prefix in Redis
	dirty  bool        // whether changes were made that need to be persisted
	fields Fieldset    // maps field names to how they should be indexed
}

const (
	cIndicesKey = "indices"
)

func (self *index_t) Conn() redis.Conn {
	return self.pool.Get()
}

func (self *index_t) Name() string {
	return self.name
}

func (self *index_t) FieldsKey() string {
	return fmt.Sprintf("fields:%s", self.Name())
}

func (self *index_t) Exists() (bool, error) {
	conn := self.pool.Get()
	defer conn.Close()

	res, err := redis.Bool(conn.Do("SISMEMBER", cIndicesKey, self.name))
	if err != nil {
		return false, errgo.Mask(err)
	}

	return res, nil
}

func (self *index_t) Load() error {
	exists, err := self.Exists()
	if err != nil {
		return errgo.Mask(err)
	}
	if !exists {
		return errgo.Newf("index '%s' does not exist", self.name)
	}

	conn := self.pool.Get()
	defer conn.Close()

	val, err := redis.Strings(conn.Do("HGETALL", self.FieldsKey()))
	if err != nil {
		return errgo.Mask(err)
	}

	self.fields = make(Fieldset)
	for k := 0; k < len(val); k += 2 {
		self.AddField(val[k], FieldType(val[k+1]))
	}

	return nil
}

// TODO: make this transactional, using version numbers/UUID and
// a LUA script
func (self *index_t) Save() error {
	conn := self.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", "indices", self.name)
	if err != nil {
		return errgo.Mask(err)
	}

	for _, field := range self.Fields() {
		fmt.Println("saving field", field.Name())
		err := field.Save()
		if err != nil {
			return errgo.Mask(err)
		}
	}

	return nil
}

func (self *index_t) Destroy() error {
	return errgo.New("index_t#Destroy not implemented")
}

func (self *index_t) Fields() Fieldset {
	return self.fields
}

func (self *index_t) AddField(name string, ty FieldType) error {
	field, err := newField(self, name, ty)
	if err != nil {
		return errgo.Mask(err)
	}
	self.Fields()[name] = field
	return nil
}
