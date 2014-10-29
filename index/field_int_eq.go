package index

import (
	"github.com/juju/errgo"
)

type field_int_eq_t struct {
	field_with_int_value_t
	field_t
}

func (self *field_int_eq_t) Save() error {
	if err := self.saveCommon(); err != nil {
		return errgo.Mask(err)
	}

	key := self.DataKey()
	keyType, err := self.keyType(key)
	if err != nil {
		return errgo.Mask(err)
	}

	switch keyType {
	case "none":
		// will be set lazily
	case "set":
		// looks good
	default:
		return errgo.Newf("key '%s' of type '%s', expected 'none' or 'set'", key, keyType)
	}

	return nil
}
