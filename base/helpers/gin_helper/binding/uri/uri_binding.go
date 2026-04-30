package uri

import (
	cBinding "core-ticket/base/helpers/gin_helper/binding"
	"core-ticket/base/helpers/gin_helper/binding/utility"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin/binding"
)

type Binding struct {
}

func (c Binding) Name() string {
	return "uriBinding"
}

func (c Binding) Bind(m map[string][]string, obj any) error {
	if err := binding.MapFormWithTag(obj, m, cBinding.Uri.RulesTag()); err != nil {
		var numError *strconv.NumError
		switch {
		case errors.As(err, &numError),
			errors.Is(err, strconv.ErrRange),
			errors.Is(err, strconv.ErrSyntax):

			return utility.ValidateMap(
				utility.ConvertUrlValuesToMapStringInterface(m),
				utility.GetRulesFromStruct(obj, cBinding.Uri),
			)
		}
		return err
	}

	return utility.ValidateStruct(obj)
}
