package crud_handler

import (
	"core-ticket/base/helpers/base_helper"
	"core-ticket/base/helpers/context_helper"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/http_helper"
	"core-ticket/base/helpers/struct_validator"
	"core-ticket/constants"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHandlerWithPathParamFunc[T interface{}, U interface{}](fn func(identity base_helper.Identity, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p T
		identity, err := context_helper.GetIdentity(c)
		if err != nil {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, err.Error(), nil)
			return
		}

		if err := c.ShouldBindUri(&p); err != nil {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, err.Error(), nil)
			return
		}

		res, handlerErr := fn(identity, p)

		if handlerErr != nil {
			if handlerErr.StatusCode == http.StatusInternalServerError {
				http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else if handlerErr.StatusCode == http.StatusBadRequest {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else {
				http_helper.ErrorResponse(c, handlerErr)
			}
			return
		}
		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func GetHandlerFunc[T interface{}, U interface{}](fn func(serverCode string, dbName string, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		if err := c.ShouldBindQuery(&d); err != nil {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, err.Error(), nil)
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		res, handlerErr := fn(serverCode, dbName, d)
		if handlerErr != nil {
			fmt.Printf("handlerErr: %+v\n", handlerErr)
			if handlerErr.StatusCode != http.StatusInternalServerError {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
				return
			}
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func GetPaginationHandlerFunc[T interface{}, U interface{}](fn func(serverCode string, dbName string, t T) ([]U, int, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		if err := c.ShouldBindQuery(&d); err != nil {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, err.Error(), nil)
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		res, count, handlerErr := fn(serverCode, dbName, d)
		if handlerErr != nil {
			if handlerErr.StatusCode != http.StatusInternalServerError {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
				return
			}
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessPaginationResponse[U](c, constants.EC_SUCCESS, "OK", res, count, page, limit)
	}
}

func FilterHandlerFunc[T interface{}](fn func(T) (http_helper.PaginationResult[T], *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		res, handlerErr := fn(d)
		if handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func CreateHandlerFunc[T interface{}, U interface{}](fn func(serverCode string, dbName string, username string, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseJSON(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			return
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		res, handlerErr := fn(serverCode, dbName, username, d)
		if handlerErr != nil {
			if handlerErr.StatusCode != http.StatusInternalServerError {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
				return
			}
			http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func UpdateHandlerFunc[T interface{}, U interface{}](fn func(serverCode string, dbName string, username string, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		_ = validate.ParseUri(c, &d)
		if handlerErr := validate.ParseJSON(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			return
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		res, handlerErr := fn(serverCode, dbName, username, d)
		if handlerErr != nil {
			if handlerErr.StatusCode != http.StatusInternalServerError {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
				return
			}
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func DeleteHandlerFunc[T interface{}, U interface{}](fn func(serverCode string, dbName string, username string, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		res, handlerErr := fn(serverCode, dbName, username, d)
		if handlerErr != nil {
			if handlerErr.StatusCode != http.StatusInternalServerError {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
				return
			}
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func MainHandlerFuncWithRequest[T interface{}, U interface{}](fn func(username string, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseJSON(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		res, handlerErr := fn(username, d)
		if handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func MainHandlerFunc[U interface{}](fn func(username string) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		res, handlerErr := fn(username)
		if handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func GetHandlerFuncWithIdentity[T interface{}, U interface{}](fn func(identity base_helper.Identity, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		if err := c.ShouldBindQuery(&d); err != nil {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, err.Error(), nil)
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		companyID, ok := context_helper.GetFloat(c, constants.CompanyID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company", nil)
			return
		}

		companyCode, ok := context_helper.GetString(c, constants.CompanyCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company Code", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		userRoleID, ok := context_helper.GetFloat(c, constants.UserRoleID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve User Role", nil)
			return
		}

		res, handlerErr := fn(base_helper.Identity{
			ServerCode:  serverCode,
			DbName:      dbName,
			CompanyID:   int(companyID),
			CompanyCode: companyCode,
			Username:    username,
			UserRoleID:  int(userRoleID),
		}, d)
		if handlerErr != nil {
			if handlerErr.StatusCode == http.StatusInternalServerError {
				http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else if handlerErr.StatusCode == http.StatusBadRequest {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else {
				http_helper.ErrorResponse(c, handlerErr)
			}
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func GetPaginationHandlerFuncWithIdentity[T interface{}, U interface{}](fn func(identity base_helper.Identity, t T) ([]U, int, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		if err := c.ShouldBindQuery(&d); err != nil {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, err.Error(), nil)
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		companyID, ok := context_helper.GetFloat(c, constants.CompanyID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company", nil)
			return
		}

		companyCode, ok := context_helper.GetString(c, constants.CompanyCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company Code", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		userRoleID, ok := context_helper.GetFloat(c, constants.UserRoleID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve User Role", nil)
			return
		}

		res, count, handlerErr := fn(base_helper.Identity{
			ServerCode:  serverCode,
			DbName:      dbName,
			CompanyID:   int(companyID),
			CompanyCode: companyCode,
			Username:    username,
			UserRoleID:  int(userRoleID),
		}, d)
		if handlerErr != nil {
			if handlerErr.StatusCode == http.StatusInternalServerError {
				http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else if handlerErr.StatusCode == http.StatusBadRequest {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else {
				http_helper.ErrorResponse(c, handlerErr)
			}
			return
		}

		http_helper.SuccessPaginationResponse[U](c, constants.EC_SUCCESS, "OK", res, count, page, limit)
	}
}

func CreateHandlerFuncWithIdentity[T interface{}, U interface{}](fn func(identity base_helper.Identity, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseJSON(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			return
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		companyID, ok := context_helper.GetFloat(c, constants.CompanyID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company", nil)
			return
		}

		companyCode, ok := context_helper.GetString(c, constants.CompanyCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company Code", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		userRoleID, ok := context_helper.GetFloat(c, constants.UserRoleID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve User Role", nil)
			return
		}

		res, handlerErr := fn(base_helper.Identity{
			ServerCode:  serverCode,
			DbName:      dbName,
			CompanyID:   int(companyID),
			CompanyCode: companyCode,
			Username:    username,
			UserRoleID:  int(userRoleID),
		}, d)
		if handlerErr != nil {
			if handlerErr.StatusCode == http.StatusInternalServerError {
				http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else if handlerErr.StatusCode == http.StatusBadRequest {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else {
				http_helper.ErrorResponse(c, handlerErr)
			}
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func UpdateHandlerFuncWithIdentity[T interface{}, U interface{}](fn func(identity base_helper.Identity, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		_ = validate.ParseUri(c, &d)
		if handlerErr := validate.ParseJSON(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			return
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		companyID, ok := context_helper.GetFloat(c, constants.CompanyID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company", nil)
			return
		}

		companyCode, ok := context_helper.GetString(c, constants.CompanyCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company Code", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		userRoleID, ok := context_helper.GetFloat(c, constants.UserRoleID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve User Role", nil)
			return
		}

		res, handlerErr := fn(base_helper.Identity{
			ServerCode:  serverCode,
			DbName:      dbName,
			CompanyID:   int(companyID),
			CompanyCode: companyCode,
			Username:    username,
			UserRoleID:  int(userRoleID),
		}, d)
		if handlerErr != nil {
			if handlerErr.StatusCode == http.StatusInternalServerError {
				http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else if handlerErr.StatusCode == http.StatusBadRequest {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			} else {
				http_helper.ErrorResponse(c, handlerErr)
			}
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func DeleteHandlerFuncWithIdentity[T interface{}, U interface{}](fn func(identity base_helper.Identity, t T) (U, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var d T

		validate := struct_validator.New()
		if handlerErr := validate.ParseUri(c, &d); handlerErr != nil {
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		serverCode, ok := context_helper.GetString(c, constants.ServerCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Server Code", nil)
			return
		}

		dbName, ok := context_helper.GetString(c, constants.DBName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Database", nil)
			return
		}

		companyID, ok := context_helper.GetFloat(c, constants.CompanyID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company", nil)
			return
		}

		companyCode, ok := context_helper.GetString(c, constants.CompanyCode)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Company Code", nil)
			return
		}

		username, ok := context_helper.GetString(c, constants.UserName)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve Username", nil)
			return
		}

		userRoleID, ok := context_helper.GetFloat(c, constants.UserRoleID)
		if !ok {
			http_helper.BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Failed to retrieve User Role", nil)
			return
		}

		res, handlerErr := fn(base_helper.Identity{
			ServerCode:  serverCode,
			DbName:      dbName,
			CompanyID:   int(companyID),
			CompanyCode: companyCode,
			Username:    username,
			UserRoleID:  int(userRoleID),
		}, d)
		if handlerErr != nil {
			switch handlerErr.StatusCode {
			case http.StatusInternalServerError:
				http_helper.ServerErrorResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			case http.StatusBadRequest:
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
			default:
				http_helper.ErrorResponse(c, handlerErr)
			}
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "", res)
	}
}

func GetAllHandlerFunc[T interface{}](fn func(identity base_helper.Identity) ([]T, *http_helper.HandlerError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, err := context_helper.GetIdentity(c)
		if err != nil {
			http_helper.HttpErrorResponse(c, error_helper.New(err))
			return
		}

		res, handlerErr := fn(identity)
		if handlerErr != nil {
			if handlerErr.StatusCode != http.StatusInternalServerError {
				http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, handlerErr.Data)
				return
			}
			http_helper.BadRequestResponse(c, handlerErr.RespCode, handlerErr.Message, nil)
			return
		}

		http_helper.SuccessResponse(c, constants.EC_SUCCESS, "OK", res)
	}
}
