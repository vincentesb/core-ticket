package form

import (
	cBinding "core-ticket/base/helpers/gin_helper/binding"
	"core-ticket/base/helpers/gin_helper/binding/utility"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin/binding"
)

const defaultMemory = 32 << 20

type Binding struct{}

func (c Binding) Name() string {
	return "formBinding"
}

func (c Binding) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}

	if err := binding.MapFormWithTag(obj, req.Form, cBinding.Form.RulesTag()); err != nil {
		var numError *strconv.NumError
		switch {
		case errors.As(err, &numError),
			errors.Is(err, strconv.ErrRange),
			errors.Is(err, strconv.ErrSyntax):
			return utility.ValidateMap(
				utility.ConvertUrlValuesToMapStringInterface(req.Form),
				utility.GetRulesFromStruct(obj, cBinding.Form),
			)
		}
		return err
	}

	return utility.ValidateStruct(obj)
}
