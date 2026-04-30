package base_helper

import (
	"encoding/json"
)

const (
	TYPE_INTEGER = "base_helper.Integer"
	TYPE_BOOLEAN = "base_helper.Boolean"
	TYPE_FLOAT   = "base_helper.Float"
)

type Integer int

func (r *Integer) UnmarshalJSON(data []byte) error {
	var property interface{}
	if err := json.Unmarshal(data, &property); err != nil {
		return err
	}

	*r = Integer(ConvertToInteger(property))
	return nil
}

type Boolean bool

func (r *Boolean) UnmarshalJSON(data []byte) error {
	var property interface{}
	if err := json.Unmarshal(data, &property); err != nil {
		return err
	}

	boolean, err := ConvertToBoolean(property)
	if err != nil {
		return err
	}
	*r = Boolean(boolean)
	return nil
}

type Float float64

func (r *Float) UnmarshalJSON(data []byte) error {
	var property interface{}
	if err := json.Unmarshal(data, &property); err != nil {
		return err
	}

	float := ConvertToFloat(property)

	*r = Float(float)
	return nil
}

type Identity struct {
	ServerCode  string
	DbName      string
	CompanyID   int
	CompanyCode string
	Username    string
	UserRoleID  int
}
