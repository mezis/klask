package index

import (
	"encoding/json"
	"github.com/juju/errgo"
)

type index_presenter_t struct {
	ID     string              `json:"id"`
	Fields []field_presenter_t `json:"fields"`
}

type field_presenter_t struct {
	Name string    `json:"name"`
	Type FieldType `json:"type"`
}

// Satisfiy the json.Marhshaler interface
func (self *index_t) MarshalJSON() ([]byte, error) {
	var presenter index_presenter_t

	presenter.ID = self.name
	presenter.Fields = make([]field_presenter_t, len(self.fields))
	k := 0
	for _, field := range self.fields {
		presenter.Fields[k].Name = field.Name()
		presenter.Fields[k].Type = field.Type()
		k++
	}

	data, err := json.Marshal(presenter)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return data, nil
}

func (self *index_t) UnmarshalJSON(data []byte) error {
	var presenter index_presenter_t

	err := json.Unmarshal(data, &presenter)
	if err != nil {
		return errgo.Mask(err)
	}

	self.name = presenter.ID
	for _, val := range presenter.Fields {
		err := self.addField(val.Name, val.Type)
		if err != nil {
			return errgo.Mask(err)
		}
	}

	return nil
}
