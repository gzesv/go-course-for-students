package homework

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Struct:
		vErrors = validateStruct(value, vErrors)
	default:
		return ErrNotStruct
	}

	if len(vErrors) == 0 {
		return nil
	}

	return vErrors
}

func validateStruct(value reflect.Value, vErrors ValidationErrors) ValidationErrors {
	if value.Kind() != reflect.Struct {
		vErrors = append(vErrors, appendError("", ErrNotStruct))
		return vErrors
	}

	for i := 0; i < value.Type().NumField(); i++ {
		valueField := value.Field(i)

		if valueField.Kind() == reflect.Struct {
			vErrors = validateStruct(valueField, vErrors)
		}

		valueTypeField := value.Type().Field(i)

		tag := valueTypeField.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}

		if !valueTypeField.IsExported() {
			vErrors = append(vErrors, appendError(valueTypeField.Name, ErrValidateForUnexportedFields))
			continue
		}

		switch valueField.Kind() {
		case reflect.Slice:
			for j := 0; j < valueField.Len(); j++ {
				ok := validateValue(valueTypeField, valueField.Index(j))
				if !ok {
					vErrors = append(vErrors, appendError(valueTypeField.Name, ErrInvalidValidatorSyntax))
					break
				}
			}
		default:
			ok := validateValue(valueTypeField, valueField)
			if !ok {
				vErrors = append(vErrors, appendError(valueTypeField.Name, ErrInvalidValidatorSyntax))
			}
		}

	}
	return vErrors
}

func validateValue(valueTypeField reflect.StructField, value reflect.Value) bool {
	tag := valueTypeField.Tag.Get("validate")
	rulesArr := strings.SplitN(tag, ":", 2)
	rName := rulesArr[0]
	rVal := rulesArr[1]

	switch rName {
	case "len":
		switch value.Kind() {
		case reflect.String:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return int64(len([]rune(value.String()))) == rVal
		}
		return false
	case "max":
		switch value.Kind() {
		case reflect.String:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return int64(len([]rune(value.String()))) <= rVal
		case reflect.Int:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return value.Int() <= rVal
		}
		return false
	case "min":
		switch value.Kind() {
		case reflect.String:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return int64(len([]rune(value.String()))) >= rVal
		case reflect.Int:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return false
			}
			return value.Int() >= rVal
		}
		return false
	case "in":
		switch value.Kind() {
		case reflect.String:
			if rVal == "" {
				return false
			}
			for _, s := range strings.Split(rVal, ",") {
				if s == value.String() {
					return true
				}
			}
			return false
		case reflect.Int:
			if rVal == "" {
				return false
			}
			for _, s := range strings.Split(rVal, ",") {
				s, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return false
				}
				if s == value.Int() {
					return true
				}
			}
			return false
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
