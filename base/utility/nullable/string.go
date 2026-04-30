package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

type String struct {
	sql.NullString
}

func StringFrom(s string) String {
	return NewString(s, true)
}

func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

func (s String) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

func (s *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return err
	}

	s.Valid = true
	return nil
}

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
}

func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

func (s String) IsZero() bool {
	return !s.Valid
}

func (s String) Equal(other String) bool {
	return s.Valid == other.Valid && (!s.Valid || s.String == other.String)
}
