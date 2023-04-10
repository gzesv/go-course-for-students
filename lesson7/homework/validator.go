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
		vErrors = validateStruct(value)
	default:
		return ErrNotStruct
	}

	if len(vErrors) == 0 {
		return nil
	}

	return vErrors
}

func validateStruct(value reflect.Value) ValidationErrors {
	var vErrors ValidationErrors

	if value.Kind() != reflect.Struct {
		vErrors = append(vErrors, appendError("", ErrNotStruct))
		return vErrors
	}

	for i := 0; i < value.Type().NumField(); i++ {
		valueField := value.Field(i)

		if valueField.Kind() == reflect.Struct {
			vErrors = validateStruct(valueField)
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
				err := validateValue(valueTypeField, valueField.Index(j))
				if err != nil {
					vErrors = append(vErrors, appendError(valueTypeField.Name, err))
					break
				}
			}
		default:
			err := validateValue(valueTypeField, valueField)
			if err != nil {
				vErrors = append(vErrors, appendError(valueTypeField.Name, err))
			}
		}

	}
	return vErrors
}

func validateValue(valueTypeField reflect.StructField, value reflect.Value) error {
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
				return ErrInvalidValidatorSyntax
			}
			if int64(len([]rune(value.String()))) == rVal {
				return nil
			}
			return errors.New("invalid string length")
		}
		return errors.New("invalid type of field")
	case "max":
		switch value.Kind() {
		case reflect.String:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			if int64(len([]rune(value.String()))) <= rVal {
				return nil
			}
			return errors.New("string len greater max")
		case reflect.Int:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			if value.Int() <= rVal {
				return nil
			}
			return errors.New("int value greater max")
		}
		return errors.New("invalid type of field")
	case "min":
		switch value.Kind() {
		case reflect.String:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			if int64(len([]rune(value.String()))) >= rVal {
				return nil
			}
			return errors.New("string length less min")
		case reflect.Int:
			rVal, err := strconv.ParseInt(rVal, 10, 64)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			if value.Int() >= rVal {
				return nil
			}
			return errors.New("int value less min")
		}
		return errors.New("invalid type of field")
	case "in":
		switch value.Kind() {
		case reflect.String:
			if rVal == "" {
				return ErrInvalidValidatorSyntax
			}
			for _, s := range strings.Split(rVal, ",") {
				if s == value.String() {
					return nil
				}
			}
			return errors.New("string value in not list")
		case reflect.Int:
			if rVal == "" {
				return ErrInvalidValidatorSyntax
			}
			for _, s := range strings.Split(rVal, ",") {
				s, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return ErrInvalidValidatorSyntax
				}
				if s == value.Int() {
					return nil
				}
			}
			return errors.New("int value in not list")
		}
		return errors.New("invalid type of field")
	}
	return ErrInvalidValidatorSyntax
}

func appendError(name string, err error) ValidationError {
	return ValidationError{
		Name: name,
		Err:  err,
	}
}
