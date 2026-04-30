package error_code

type ValidationErrorCode ErrorCode

const (
	IsRequired       ValidationErrorCode = "00023"
	MinLength        ValidationErrorCode = "00024"
	MaxLength        ValidationErrorCode = "00025"
	IsNumeric        ValidationErrorCode = "00026"
	GreaterThan      ValidationErrorCode = "00027"
	LessThan         ValidationErrorCode = "00028"
	IsUnique         ValidationErrorCode = "00029"
	IsBoolean        ValidationErrorCode = "00030"
	Email            ValidationErrorCode = "00031"
	PhoneNumber      ValidationErrorCode = "00032"
	Url              ValidationErrorCode = "00033"
	Integer          ValidationErrorCode = "00034"
	Float            ValidationErrorCode = "00035"
	Date             ValidationErrorCode = "00036"
	In               ValidationErrorCode = "00037"
	NotExistsInDB    ValidationErrorCode = "00038"
	GreaterThanEqual ValidationErrorCode = "00039"
	LessThanEqual    ValidationErrorCode = "00040"
	Alpha            ValidationErrorCode = "00041"
	AlphaNum         ValidationErrorCode = "00042"
	AlphaNumSpace    ValidationErrorCode = "00043"
	Regex            ValidationErrorCode = "00044"
	Custom           ValidationErrorCode = "00045"
	NotEmoji         ValidationErrorCode = "00046"
	Decimal          ValidationErrorCode = "00047"
	NotEqual         ValidationErrorCode = "00048"
)

func (vec ValidationErrorCode) String() string {
	return EC_ + string(vec)
}
