package rpc

import (
	"reflect"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
	"github.com/pkg/errors"
)

func Marshal(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	var response []byte

	reflectValue := reflect.ValueOf(v)
	reflectElem := reflectValue.Elem()

	if reflectElem.Kind() != reflect.Struct {
		logs.Error("value passed in is not of structure type")
		return nil, custom_errors.NewMarshallerError(errors.Errorf("value passed in is not of structure type"))
	}

	for i := 0; i < reflectElem.NumField(); i++ {
		field := reflectElem.FieldByName(reflectElem.Type().Field(i).Name)
		if field.IsValid() {
			fieldKind := reflectElem.Type().Field(i).Type.Kind()
			if fieldKind == reflect.Interface {
				if field.IsNil() {
					continue
				}
				fieldKind = reflect.TypeOf(field.Interface()).Elem().Kind()
			}
			switch fieldKind {
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint8, reflect.Float64, reflect.String:
				err := marshalPrimitive(&response, fieldKind, field)
				if err != nil {
					return nil, err
				}
			case reflect.Slice:
				err := marshalArray(&response, field, reflectElem.Type().Field(i).Type.Elem().Kind())
				if err != nil {
					return nil, err
				}
			case reflect.Struct:
				err := marshalStruct(&response, field)
				if err != nil {
					return nil, err
				}
			default:
				logs.Error("unimplemented type: %v", fieldKind)
				return nil, custom_errors.NewMarshallerError(errors.Errorf("unimplemented type, type: %v", fieldKind))
			}
		}
	}
	return response, nil
}

func marshalPrimitive(response *[]byte, fieldKind reflect.Kind, field reflect.Value) error {
	switch fieldKind {
	case reflect.Int64:
		*response = append(*response, bytes.Int64ToBytes(field.Interface().(int64))...)
	case reflect.Int:
		*response = append(*response, bytes.Int64ToBytes(int64(field.Interface().(int)))...)
	case reflect.Int32:
		*response = append(*response, bytes.Int32ToBytes(field.Interface().(int32))...)
	case reflect.Uint8:
		*response = append(*response, field.Convert(reflect.ValueOf(uint8(1)).Type()).Interface().(uint8))
	case reflect.Float64:
		*response = append(*response, bytes.Float64ToBytes(field.Interface().(float64))...)
	case reflect.String:
		*response = append(*response, []byte(field.Interface().(string))...)
		*response = append(*response, stringTerminator)
	default:
		logs.Error("unimplemented type: %v", fieldKind)
		return custom_errors.NewMarshallerError(errors.Errorf("unimplemented type"))
	}
	return nil
}

func marshalArray(response *[]byte, field reflect.Value, elementType reflect.Kind) error {
	sizeOfSlice := field.Len()
	*response = append(*response, bytes.Int64ToBytes(int64(sizeOfSlice))...)

	switch elementType {
	case reflect.Int:
		slice := field.Interface().([]int)
		for _, v := range slice {
			*response = append(*response, bytes.Int64ToBytes(int64(v))...)
		}
	case reflect.Int64:
		slice := field.Interface().([]int64)
		for _, v := range slice {
			*response = append(*response, bytes.Int64ToBytes(v)...)
		}
	case reflect.Int32:
		slice := field.Interface().([]int32)
		for _, v := range slice {
			*response = append(*response, bytes.Int32ToBytes(v)...)
		}
	case reflect.Uint8:
		slice := field.Interface().([]uint8)
		for _, v := range slice {
			*response = append(*response, v)
		}
	case reflect.Float64:
		slice := field.Interface().([]float64)
		for _, v := range slice {
			*response = append(*response, bytes.Float64ToBytes(v)...)
		}
	case reflect.String:
		slice := field.Interface().([]string)
		for _, v := range slice {
			*response = append(*response, []byte(v)...)
			*response = append(*response, stringTerminator)
		}
	case reflect.Struct:
		for i := 0; i < sizeOfSlice; i++ {
			val := field.Index(i)
			err := marshalStruct(response, val)
			if err != nil {
				return err
			}
		}
	default:
		logs.Error("unimplemented type: %v, name: %v", elementType, field.Type().Field(0).Name)
		return custom_errors.NewMarshallerError(errors.Errorf("unimplemented type"))

	}
	return nil
}

func marshalStruct(response *[]byte, reflectValue reflect.Value) error {
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.FieldByName(reflectValue.Type().Field(i).Name)
		if field.IsValid() {
			fieldKind := reflectValue.Type().Field(i).Type.Kind()
			switch fieldKind {
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint8, reflect.Float64, reflect.String:
				err := marshalPrimitive(response, fieldKind, field)
				if err != nil {
					return err
				}
			case reflect.Slice:
				err := marshalArray(response, field, reflectValue.Type().Field(i).Type.Elem().Kind())
				if err != nil {
					return err
				}
			case reflect.Struct:
				err := marshalStruct(response, field)
				if err != nil {
					return err
				}
			default:
				logs.Error("unimplemented type: %v", fieldKind)
				return custom_errors.NewMarshallerError(errors.Errorf("unimplemented type, type: %v", fieldKind))
			}
		}
	}
	return nil
}
