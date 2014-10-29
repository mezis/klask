package index

import (
	"github.com/juju/errgo"
)

type Id int64

type Record interface {
	// Unique record id; negative if unsaved
	Id() Id

	// Gets the value of an attribute.
	// Errors is the field is unknown.
	Get(string) (interface{}, error)

	// Changes the value of an attribute; removes if set to nil.
	// Errors when the field is unknown or the value is of the wrong type.
	Set(string, interface{}) error

	Index() Index

	// Checks whether all fields have the correct type of value
	IsValid() bool

	Persist() error
}

type attrs_t map[string]interface{}

type record_t struct {
	id    Id
	attrs attrs_t
	idx   Index
}

func (self *index_t) New() Record {
	record := new(record_t)
	record.id = -1
	record.attrs = make(attrs_t)
	record.idx = self

	for key, _ := range self.Fields() {
		record.attrs[key] = nil
	}
	return record
}

func (self *index_t) Find(id Id) (Record, error) {
	return nil, errgo.New("not implemented")
}

func (self *record_t) Id() Id {
	return self.id
}

func (self *record_t) Get(key string) (interface{}, error) {
	value, ok := self.attrs[key]
	if !ok {
		return nil, errgo.Newf("no such field '%s'", key)
	}
	return value, nil
}

func (self *record_t) Set(key string, value interface{}) error {
	field, ok := self.idx.Fields()[key]
	if !ok {
		return errgo.Newf("no such field '%s'", key)
	}
	if err := field.CheckValidValue(value); err != nil {
		return errgo.Mask(err)
	}
	self.attrs[key] = value
	return nil
}

func (self *record_t) Index() Index {
	return self.idx
}

func (self *record_t) IsValid() bool {
	return false
}

func (self *record_t) Persist() error {
	return errgo.New("not implemented")
}
