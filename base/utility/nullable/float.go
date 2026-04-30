package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type Float struct {
	sql.NullFloat64
}

func NewFloat(f float64, valid bool) Float {
	return Float{
		NullFloat64: sql.NullFloat64{
			Float64: f,
			Valid:   valid,
		},
	}
}

func FloatFrom(f float64) Float {
	return NewFloat(f, true)
}

func FloatFromPtr(f *float64) Float {
	if f == nil {
		return NewFloat(0, false)
	}
	return NewFloat(*f, true)
}

func (f Float) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

func (f *Float) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		f.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &f.Float64); err != nil {
		return err
	}

	f.Valid = true
	return nil
}

func (f *Float) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	f.Valid = true
	return err
}

func (f Float) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	if math.IsInf(f.Float64, 0) || math.IsNaN(f.Float64) {
		return nil, &json.UnsupportedValueError{
			Value: reflect.ValueOf(f.Float64),
			Str:   strconv.FormatFloat(f.Float64, 'g', -1, 64),
		}
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

func (f Float) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(f.Float64, 'f', -1, 64)), nil
}

func (f *Float) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
}

func (f Float) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

func (f Float) IsZero() bool {
	return !f.Valid
}

func (f Float) Equal(other Float) bool {
	return f.Valid == other.Valid && (!f.Valid || f.Float64 == other.Float64)
}
