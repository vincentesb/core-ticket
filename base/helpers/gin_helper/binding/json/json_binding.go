package json

import (
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/gin_helper/binding/utility"
	"core-ticket/constants/error_code"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Binding struct{}

func (c Binding) Name() string {
	return "jsonBinding"
}

func (c Binding) Bind(req *http.Request, obj any) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}

	if err := c.decodeJSON(req.Body, obj); err != nil {
		var jsonErr *json.UnmarshalTypeError
		if errors.As(err, &jsonErr) {
			return error_helper.New(error_helper.NewValidationError(
				error_code.ValidationErrorCode(error_code.ValidationError),
				"custom",
				jsonErr.Field,
				"",
				0,
				fmt.Sprintf("%s invalid type sent (%s) to %s", jsonErr.Field, jsonErr.Value, jsonErr.Type.String()),
			), error_code.ValidationError)
		}
		return error_helper.New(errors.New("invalid JSON syntax please check your payload"), error_code.ValidationError)
	}

	return utility.ValidateStruct(obj)
}

func (c Binding) decodeJSON(r io.Reader, obj any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
