package main

import (
	"reflect"
	"errors"
)

var (
	ErrInvalidType = errors.New("invalid type")
)

func i2s(data interface{}, out interface{}) error {
	var err error
	switch data.(type) {
	case []interface{}:
		val := reflect.Indirect(reflect.ValueOf(out))
		typePtr := reflect.TypeOf(out).Elem()
		if typePtr.Kind() != reflect.Slice {
			return ErrInvalidType
		}
		refType := typePtr.Elem()
		slice := reflect.Zero(reflect.SliceOf(refType))
		for _, item := range data.([]interface{}) {
			elem := reflect.Indirect(reflect.New(refType))
			if err = unpackReflect(item, elem.Addr().Interface()); err != nil {
				break
			}
			slice = reflect.Append(slice, elem)
		}
		val.Set(slice)
	case interface{}:
		err = unpackReflect(data, out)
	default:
		err = ErrInvalidType
	}
	return err
}

func unpackReflect(data interface{}, out interface{}) error {
	typePtr := reflect.ValueOf(out)
	if typePtr.Kind() != reflect.Ptr {
		return ErrInvalidType
	}
	val := reflect.Indirect(reflect.ValueOf(out))
	m := data.(map[string]interface{})
	for k, v := range m {
		valueField := val.FieldByName(k)
		ref := reflect.ValueOf(v).Kind()
		if err := checkTypeCompatibility(ref, valueField.Kind()); err != nil {
			return err
		}
		switch v.(type) {
		case float64:
			valueField.Set(reflect.ValueOf(int(v.(float64))))
		case string:
			valueField.Set(reflect.ValueOf(v.(string)))
		case bool:
			valueField.Set(reflect.ValueOf(v.(bool)))
		case map[string]interface{}:
			i2s(v, valueField.Addr().Interface())
		case []interface{}:
			s := reflect.ValueOf(v)
			refType := reflect.TypeOf(valueField.Interface()).Elem()
			slice := reflect.Zero(reflect.SliceOf(refType))
			for i := 0; i < s.Len(); i++ {
				elem := reflect.Indirect(reflect.New(refType))
				i2s(s.Index(i).Interface(), elem.Addr().Interface())
				slice = reflect.Append(slice, elem)
			}
			valueField.Set(slice)
		default:
			return ErrInvalidType
		}
	}
	return nil
}

func checkTypeCompatibility(lhs, rhs reflect.Kind) error {
	if lhs == reflect.Float64 {
		lhs = reflect.Int
	}
	if lhs == reflect.Map {
		lhs = reflect.Struct
	}
	if lhs != rhs{
		return ErrInvalidType
	}
	return nil
}
