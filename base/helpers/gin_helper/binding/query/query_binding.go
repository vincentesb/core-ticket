package query

import (
	cBinding "core-ticket/base/helpers/gin_helper/binding"
	"core-ticket/base/helpers/gin_helper/binding/utility"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
)

type Binding struct {
}

func (c Binding) Name() string {
	return "queryBinding"
}

func parseQueryString(rawQuery string) (url.Values, error) {
	values := make(url.Values)
	if rawQuery == "" {
		return values, nil
	}

	for _, param := range strings.Split(rawQuery, "&") {
		if param == "" {
			continue
		}

		eqIndex := strings.Index(param, "=")
		if eqIndex == -1 {
			key, err := url.QueryUnescape(param)
			if err != nil {
				return nil, err
			}
			values[key] = append(values[key], "")
		} else {
			key, err := url.QueryUnescape(param[:eqIndex])
			if err != nil {
				return nil, err
			}
			value, err := url.QueryUnescape(param[eqIndex+1:])
			if err != nil {
				return nil, err
			}
			values[key] = append(values[key], value)
		}
	}

	return values, nil
}

func (c Binding) Bind(req *http.Request, obj any) error {
	values, err := parseQueryString(req.URL.RawQuery)
	if err != nil {
		values = req.URL.Query()
	}

	if err := binding.MapFormWithTag(obj, values, cBinding.Query.RulesTag()); err != nil {
		var numError *strconv.NumError
		switch {
		case errors.As(err, &numError),
			errors.Is(err, strconv.ErrRange),
			errors.Is(err, strconv.ErrSyntax):

			return utility.ValidateMap(
				utility.ConvertUrlValuesToMapStringInterface(values),
				utility.GetRulesFromStruct(obj, cBinding.Query),
			)
		}
		return err
	}

	return utility.ValidateStruct(obj)
}
