package index

import (
	"encoding/json"
	"fmt"
	"github.com/juju/errgo"
)

type record_presenter_t map[string]interface{}

func (self *record_t) MarshalJSON() ([]byte, error) {
	var (
		// presenter record_presenter_t
		err error
	)

	presenter := make(record_presenter_t)

	presenter["id"] = self.Id()

	for key, _ := range self.Index().Fields() {
		value, err := self.Get(key)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		presenter[key] = &value
	}

	data, err := json.Marshal(presenter)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return data, nil
}

func (self *record_t) UnmarshalJSON(data []byte) error {
	var (
		presenter record_presenter_t
		err       error
	)

	err = json.Unmarshal(data, &presenter)
	if err != nil {
		return errgo.Mask(err)
	}

	for key, value := range presenter {
		if key == "id" {
			switch value_ := value.(type) {
			case float64:
				self.id = Id(value_)
			case string:
				_, err := fmt.Sscanf(value_, "%d", &self.id)
				if err != nil {
					return errgo.Mask(err)
				}
			default:
				return errgo.Newf("id field of unexpected type")
			}
			continue
		}
		err = self.Set(key, value)
		if err != nil {
			fmt.Printf("error setting %s to %v (%T)\n", key, value, value)
			return errgo.Mask(err)
		}
	}

	return nil
}
