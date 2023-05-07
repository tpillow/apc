package apcgen

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tpillow/apc/pkg/apc"
)

func valueSetFieldOrAppendKind(rawVal any, valKind reflect.Kind, field reflect.Value) {
	maybeAppendValToSliceTrueIfNot := func(val any) bool {
		if field.Kind() != reflect.Slice {
			return true
		}
		field.Set(reflect.Append(field, reflect.ValueOf(val)))
		return false
	}

	panicUnsettable := func(val any, exp string) {
		panic(fmt.Sprintf("cannot set field to value '%v': cannot convert %v", val, exp))
	}

	switch valKind {
	case reflect.String:
		switch val := rawVal.(type) {
		case string:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetString(val)
			}
		default:
			panicUnsettable(rawVal, "to string")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch val := rawVal.(type) {
		case int, int8, int16, int32, int64:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetInt(val.(int64))
			}
		case string:
			if cVal, err := strconv.ParseInt(val, 10, 64); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetInt(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to int from string")
			}
		default:
			panicUnsettable(rawVal, "to int")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch val := rawVal.(type) {
		case uint, uint8, uint16, uint32, uint64:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetUint(val.(uint64))
			}
		case string:
			if cVal, err := strconv.ParseUint(val, 10, 64); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetUint(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to uint from string")
			}
		default:
			panicUnsettable(rawVal, "to uint")
		}
	case reflect.Float32, reflect.Float64:
		switch val := rawVal.(type) {
		case float32, float64:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetFloat(val.(float64))
			}
		case string:
			if cVal, err := strconv.ParseFloat(val, 64); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetFloat(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to float from string")
			}
		default:
			panicUnsettable(rawVal, "to float")
		}
	case reflect.Bool:
		switch val := rawVal.(type) {
		case bool:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetBool(val)
			}
		case string:
			if cVal, err := strconv.ParseBool(val); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetBool(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to bool from string")
			}
		case apc.MaybeValue[any]:
			field.SetBool(!val.IsNil())
		case apc.MaybeValue[string]:
			field.SetBool(!val.IsNil())
		case apc.MaybeValue[apc.Token]:
			field.SetBool(!val.IsNil())
		default:
			panicUnsettable(rawVal, "to bool")
		}
	case reflect.Pointer:
		switch val := rawVal.(type) {
		case apc.MaybeValue[any]:
			if !val.IsNil() {
				if maybeAppendValToSliceTrueIfNot(val.Value()) {
					field.Set(reflect.ValueOf(val.Value()))
				}
			}
		default:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.Set(reflect.ValueOf(val))
			}
		}
	default:
		panicUnsettable(rawVal, fmt.Sprintf("unsupported value kind %v", valKind))
	}
}
