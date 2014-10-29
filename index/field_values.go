package index

import (
	"github.com/juju/errgo"
)

type field_with_int_value_t struct {
}

func (self *field_with_int_value_t) CheckValidValue(value interface{}) error {
	switch value.(type) {
	case int:
		return nil
	default:
		return errgo.Newf("bad value of type '%T'; expected 'int'", value)
	}
	return nil
}
