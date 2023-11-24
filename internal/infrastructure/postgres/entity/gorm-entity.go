package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Entity interface {
	GetID() interface{}
}

type JSONData map[string]interface{}

func (a JSONData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
