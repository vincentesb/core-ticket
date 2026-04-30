package binding

type Type string

const (
	JSON  Type = "json"
	Form  Type = "form"
	Uri   Type = "uri"
	Query Type = "query"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) IsValid() bool {
	switch t {
	case JSON, Form, Uri, Query:
		return true
	default:
		return false
	}
}

func (t Type) RulesTag() string {
	switch t {
	case JSON:
		return "json"
	case Form:
		return "form"
	case Uri:
		return "uri"
	case Query:
		return "form"
	default:
		return "form"
	}
}
