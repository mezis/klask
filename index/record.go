package index

import (
	"encoding/json"
	"github.com/juju/errgo"
)

type Id int64
type IdList []Id

type Record interface {
	json.Marshaler
	json.Unmarshaler

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
	// FIXME: not implemented or used.
	// IsValid() bool

	// Persists the record.
	Persist() error

	// Destroys this record. Erros if the records wasn't saved beforehand.
	Destroy() error
}

var (
	ErrorNotFound = errgo.New("record not found")
	ErrorUnsaved  = errgo.New("record not saved")
)

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
	if value, err := field.CheckValidValue(value); err != nil {
		return errgo.Mask(err)
	} else {
		self.attrs[key] = value
		return nil
	}
}

func (self *record_t) Index() Index {
	return self.idx
}



func (self *record_t) Persist() error {
	var err error

	for name, field := range self.idx.Fields() {
		value, ok := self.attrs[name]
		if !ok {
			continue
		}
		if err = field.Add(self.id, value); err != nil {
			return errgo.Mask(err)
		}
	}
	if err = self.idx.Persist(self); err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (self *record_t) Destroy() error {
	if self.id < 0 {
		return errgo.Mask(ErrorUnsaved)
	}
	// TODO: ignore non-connection errors?
	if err := self.idx.Del(self.id); err != nil {
		return errgo.Mask(err)
	}
	return nil
}
