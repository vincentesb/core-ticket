package context_helper

import (
	"core-ticket/base/helpers/base_helper"
	"core-ticket/constants"
	"errors"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

func SetBulk(c *gin.Context, v map[string]interface{}) {
	for s, a := range v {
		c.Set(s, a)
	}
}

func GetHeader(c *gin.Context, key string) string {
	return c.Request.Header.Get(key)
}

func GetAuth(c *gin.Context) string {
	return GetHeader(c, "Authorization")
}

func GetHeaders(c *gin.Context, keys []string) map[string]string {
	var headers = map[string]string{}
	for _, key := range keys {
		headers[key] = GetHeader(c, key)
	}
	return headers
}

func GetString(c *gin.Context, key string) (string, bool) {
	res, exist := c.Get(key)
	val, ok := res.(string)
	if !ok || !exist {
		exist = false
	}
	return val, exist
}

func GetIdentity(c *gin.Context) (base_helper.Identity, error) {
	var (
		username    = ""
		serverCode  = ""
		dbName      = ""
		companyCode = ""
		companyID   = 0
		userRoleID  = 0
	)

	if val, ok := GetString(c, constants.UserName); ok {
		username = val
	} else {
		return base_helper.Identity{}, errors.New("failed to retrieve Username")
	}

	if val, ok := GetString(c, constants.DBName); ok {
		dbName = val
	} else {
		return base_helper.Identity{}, errors.New("failed to retrieve DB Name")
	}

	if val, ok := GetString(c, constants.ServerCode); ok {
		serverCode = val
	} else {
		return base_helper.Identity{}, errors.New("failed to retrieve Server Code")
	}

	if val, ok := GetString(c, constants.CompanyCode); ok {
		companyCode = val
	} else {
		return base_helper.Identity{}, errors.New("failed to retrieve Company Code")
	}

	if val, ok := GetFloat(c, constants.CompanyID); ok {
		companyID = int(val)
	} else {
		if val, ok := GetInt(c, constants.CompanyID); ok {
			companyID = val
		} else {
			return base_helper.Identity{}, errors.New("failed to retrieve Company ID")
		}
	}

	if val, ok := GetFloat(c, constants.UserRoleID); ok {
		userRoleID = int(val)
	} else {
		if val, ok := GetInt(c, constants.UserRoleID); ok {
			userRoleID = val
		} else {
			return base_helper.Identity{}, errors.New("failed to retrieve User Role ID")
		}
	}

	return base_helper.Identity{
		Username:    username,
		ServerCode:  serverCode,
		DbName:      dbName,
		CompanyCode: companyCode,
		CompanyID:   companyID,
		UserRoleID:  userRoleID,
	}, nil
}

func GetInt(c *gin.Context, key string) (int, bool) {
	res, exist := c.Get(key)
	val, ok := res.(int)
	if !ok || !exist {
		exist = false
	}
	return val, exist
}

func GetFloat(c *gin.Context, key string) (float64, bool) {
	res, exist := c.Get(key)
	val, ok := res.(float64)
	if !ok || !exist {
		exist = false
	}
	return val, exist
}

func ParamInt(c *gin.Context, key string) (int, bool) {
	res := c.Param(key)
	val, err := strconv.Atoi(res)
	var exist = true
	if err != nil {
		exist = false
	}
	return val, exist
}

func FormatDateTime(dateTime *string) *string {
	if dateTime != nil {
		parsedDateTime, _ := time.Parse(time.RFC3339, *dateTime)
		formattedDateTime := parsedDateTime.Format(time.DateTime)
		return &formattedDateTime
	}
	return nil
}

func FormatDate(date *string) *string {
	if date != nil {
		parsedDate, _ := time.Parse(time.RFC3339, *date)
		formattedDate := parsedDate.Format(time.DateOnly)
		return &formattedDate
	}
	return nil
}

func IsEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		if !reflect.DeepEqual(v.Field(i).Interface(), reflect.Zero(v.Field(i).Type()).Interface()) {
			return false
		}
	}

	return true
}

func GetSubtotal(qty float64, price float64, discount float64, vatValue float64, taxRate float64, vatTotal float64) float64 {
	beforeVatSubtotal := (price * qty) - discount
	if vatTotal == 0 {
		vatTotal = beforeVatSubtotal * vatValue / 100
	}

	taxTotal := beforeVatSubtotal * taxRate / 100

	return beforeVatSubtotal + vatTotal + taxTotal
}

func Contains(slice interface{}, element interface{}) bool {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice {
		panic("Contains() requires a slice as the first argument")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		current := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(current, element) {
			return true
		}
	}

	return false
}

func ShowNumericValueInWords(numericVal int) string {
	var numericWords = []string{
		"",
		"satu",
		"dua",
		"tiga",
		"empat",
		"lima",
		"enam",
		"tujuh",
		"delapan",
		"sembilan",
		"sepuluh",
		"sebelas",
	}
	wordValue := ""
	switch {
	case numericVal < 12:
		wordValue = " " + numericWords[numericVal]
	case numericVal < 20:
		wordValue = ShowNumericValueInWords(numericVal-10) + " belas"
	case numericVal < 100:
		wordValue = ShowNumericValueInWords(numericVal/10) + " puluh" + ShowNumericValueInWords(numericVal%10)
	case numericVal < 200:
		wordValue = " seratus" + ShowNumericValueInWords(numericVal-100)
	case numericVal < 1000:
		wordValue = ShowNumericValueInWords(numericVal/100) + " ratus" + ShowNumericValueInWords(numericVal%100)
	case numericVal < 2000:
		wordValue = " seribu" + ShowNumericValueInWords(numericVal-1000)
	case numericVal < 1000000:
		wordValue = ShowNumericValueInWords(numericVal/1000) + " ribu" + ShowNumericValueInWords(numericVal%1000)
	case numericVal < 1000000000:
		wordValue = ShowNumericValueInWords(numericVal/1000000) + " juta" + ShowNumericValueInWords(numericVal%1000000)
	case numericVal < 1000000000000:
		wordValue = ShowNumericValueInWords(numericVal/1000000000) + " miliar" + ShowNumericValueInWords(int(math.Mod(float64(numericVal), 1000000000)))
	case numericVal < 1000000000000000:
		wordValue = ShowNumericValueInWords(numericVal/1000000000000) + " triliun" + ShowNumericValueInWords(int(math.Mod(float64(numericVal), 1000000000000)))
	case numericVal < 1000000000000000000:
		wordValue = ShowNumericValueInWords(numericVal/1000000000000000) + " kuadriliun" + ShowNumericValueInWords(int(math.Mod(float64(numericVal), 1000000000000000)))
	}

	return wordValue
}

func Ucwords(s string) string {
	// Split the string into words.
	words := strings.Fields(s)

	// Iterate over each word and capitalize the first letter.
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	// Join the words back into a single string.
	return strings.Join(words, " ")
}
