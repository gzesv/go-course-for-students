package homework

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Name string
	Err  error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var result string
	for _, err := range v {
		result += fmt.Sprintf("%v", err.Err)
	}
	return result
}

func Validate(v any) error {
	var vErrors ValidationErrors

	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("validate")

		if len(tag) == 0 {
			continue
		}

		valueOf := val.Field(i)

		fieldName := val.Type().Field(i).Name

		if unicode.IsLower(rune(fieldName[0])) {
			vErrors = append(vErrors, appendError(fieldName, ErrValidateForUnexportedFields))
			continue
		}

		var values []interface{}
		values = getCollectValues(values, valueOf)

		if valueOf.Kind() == reflect.Slice {
			b := validateSlice(values, tag)
			if !b {
				vErrors = append(vErrors, appendError(fieldName, ErrInvalidValidatorSyntax))
			}
		} else {
			for _, value := range values {
				vErrors = append(vErrors, valIntAndStr(value, fieldName, tag)...)
			}
		}
	}

	if len(vErrors) == 0 {
		return nil
	}

	return vErrors
}

func getCollectValues(values []interface{}, value reflect.Value) []interface{} {
	switch value.Kind() {
	case reflect.Slice:
		for j := 0; j < value.Len(); j++ {
			values = append(values, value.Index(j).Interface())
		}
		return values
	default:
		values = append(values, value.Interface())
		return values
	}
}
func validateSlice(values []interface{}, tag string) bool {
	rulesArr := strings.SplitN(tag, ":", 2)
	rName := rulesArr[0]
	rVal := rulesArr[1]

	switch rName {
	case "len":
		oVal, err := strconv.ParseInt(rVal, 10, 64)
		if err != nil {
			return false
		}
		for _, v := range values {
			if len([]rune(interfaceToString(v))) != int(oVal) {
				return false
			}
		}
		return true
	case "max":
		_, err := strconv.ParseInt(rVal, 10, 64)
		if err != nil {
			return false
		}
		for _, v := range values {
			if rVal < interfaceToString(v) {
				return false
			}
		}
		return true
	case "min":
		_, err := strconv.ParseInt(rVal, 10, 64)
		if err != nil {
			return false
		}
		for _, v := range values {
			if rVal > interfaceToString(v) {
				return false
			}
		}
		return true
	case "in":
		var counter int
		if rVal == "" {
			return false
		}
		for _, v := range values {
			for _, s := range strings.Split(rVal, ",") {
				if s == interfaceToString(v) {
					counter++
					break
				}
			}
		}
		return len(values) == counter
	}
	return false
}

func valIntAndStr(value interface{}, name string, rules string) ValidationErrors {
	var vError ValidationErrors

	rulesArr := strings.SplitN(rules, ":", 2)
	rName := rulesArr[0]
	rVal := rulesArr[1]

	ok := validateValue(interfaceToString(value), rName, rVal)
	if !ok {
		vError = append(vError, appendError(name, ErrInvalidValidatorSyntax))
	}

	if len(vError) == 0 {
		return nil
	}

	return vError
}

func validateValue(val string, rName string, rVal string) bool {
	switch rName {
	case "len":
		oVal, err := strconv.ParseInt(rVal, 10, 64)
		if err != nil {
			return false
		}

		return len([]rune(val)) == int(oVal)
	case "max":
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			oVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return len(val) <= int(oVal)
		}

		ruleVal, err := strconv.ParseInt(rVal, 10, 64)
		if err != nil {
			return false
		}

		return v <= ruleVal
	case "min":
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			oVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return len(val) >= int(oVal)
		}

		ruleVal, err := strconv.ParseInt(rVal, 10, 64)
		if err != nil {
			return false
		}

		return v >= ruleVal
	case "in":
		if rVal == "" {
			return false
		}
		for _, s := range strings.Split(rVal, ",") {
			if s == val {
				return true
			}
		}
		return false
	}
	return false
}

func appendError(name string, err error) ValidationError {
	return ValidationError{
		Name: name,
		Err:  err,
	}
}

func interfaceToString(i interface{}) string {
	format := "%s"
	switch i.(type) {
	case float32, float64:
		format = "%f"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		format = "%d"
	}
	return fmt.Sprintf(format, i)
}
