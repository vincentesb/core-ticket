package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
)

type Bool struct {
	sql.NullBool
}

func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

func (b Bool) ValueOrZero() bool {
	return b.Valid && b.Bool
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		b.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &b.Bool); err != nil {
		return err
	}

	b.Valid = true
	return nil
}

func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		return errors.New("null: invalid input for UnmarshalText:" + str)
	}
	b.Valid = true
	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

func (b Bool) IsZero() bool {
	return !b.Valid
}

func (b Bool) Equal(other Bool) bool {
	return b.Valid == other.Valid && (!b.Valid || b.Bool == other.Bool)
}
