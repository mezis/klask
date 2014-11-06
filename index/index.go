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

	Field(string) (Field, error)

	// Return a new, unsaved record
	New() Record

	// Find a saved record by identifier
	// Find(id Id) (Record, error)

	// True if a record with this `id` has been saved.
	// Faster then `Find` (does not load record data)
	HasRecord(id Id) (bool, error)

	// Allocates an id if necessary, and remembers the Id for the record.
	// Does not save the record fields.
	// Should be called from Record#Persist.
	Persist(Record) error

	// Removes a persisted record by id
	Del(id Id) error

	// A Redis set holding the IDs of all known records
	RecordsKey() string

	NewTempKey() (Tempkey, error)
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

	val, err := redis.Strings(conn.Do("HGETALL", self.fieldsKey()))
	if err != nil {
		return errgo.Mask(err)
	}

	self.fields = make(Fieldset)
	for k := 0; k < len(val); k += 2 {
		self.addField(val[k], FieldType(val[k+1]))
	}

	return nil
}

// TODO: make this transactional, using version numbers/UUID and
// a LUA script
func (self *index_t) Save() error {
	var err error = nil

	conn := self.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SADD", "indices", self.name)
	if err != nil {
		return errgo.Mask(err)
	}

	for _, field := range self.Fields() {
		fmt.Println("saving field", field.Name())
		// err = field.Check()
		// if err != nil {
		// 	return errgo.Mask(err)
		// }

		err = self.saveField(field)
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

func (self *index_t) Field(name string) (Field, error) {
	field, ok := self.fields[name]
	if !ok {
		return nil, errgo.Newf("no field named '%s' in index '%s'", name, self.name)
	}
	return field, nil
}

func (self *index_t) addField(name string, ty FieldType) error {
	field, err := newField(self, name, ty)
	if err != nil {
		return errgo.Mask(err)
	}
	self.fields[name] = field
	return nil
}

func (self *index_t) saveField(field Field) error {
	conn := self.Conn()
	defer conn.Close()

	key := self.fieldsKey()

	// TODO: this needs to be transactional
	val, err := conn.Do("HGET", key, self.name)
	if err != nil {
		return errgo.Mask(err)
	}
	if val != nil {
		if val, _ := redis.String(val, err); val != string(field.Type()) {
			return errgo.Newf("field '%s' already has type '%s'", field.Name(), val)
		}
	}

	_, err = conn.Do("HSET", key, field.Name(), field.Type())
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *index_t) Persist(record Record) error {
	if record.Id() < 0 {
		return errgo.New("not implemented")
	}

	conn := self.Conn()
	defer conn.Close()

	_, err := conn.Do("SADD", self.RecordsKey(), record.Id())
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *index_t) HasRecord(id Id) (bool, error) {
	conn := self.Conn()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("SISMEMBER", self.RecordsKey(), id))
	if err != nil {
		return false, errgo.Mask(err)
	}
	return ok, nil
}

func (self *index_t) Del(id Id) error {
	var result error = nil

	for _, field := range self.fields {
		err := field.Del(id)
		if err != nil {
			result = errgo.Mask(err)
		}
	}

	conn := self.Conn()
	defer conn.Close()

	_, err := conn.Do("SREM", self.RecordsKey(), id)
	if err != nil {
		result = errgo.Mask(err)
	}

	return result
}

func (self *index_t) fieldsKey() string {
	return fmt.Sprintf("fields:%s", self.name)
}

func (self *index_t) RecordsKey() string {
	return fmt.Sprintf("records:%s", self.name)
}
