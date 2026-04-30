package struct_validator

import "github.com/gin-gonic/gin/binding"

/*
Deprecated: this type is deprecated and should not be used.
*/
type Binding string

var (
	/*
		Deprecated: this var is deprecated and should not be used.
	*/
	BindingJSON Binding = "json"
	/*
		Deprecated: this var is deprecated and should not be used.
	*/
	BindingQuery Binding = "form"
	/*
		Deprecated: this var is deprecated and should not be used.
	*/
	BindingUri Binding = "uri"
)

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) stringToHttpBinding(b Binding) binding.Binding {
	switch b {
	case BindingJSON:
		return binding.JSON
	case BindingQuery:
		return binding.Query
	default:
		return binding.JSON
	}
}
