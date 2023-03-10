package rpc

import (
	"reflect"

	"github.com/cyiafn/flight_information_system/server/custom_errors"
	"github.com/cyiafn/flight_information_system/server/logs"
	"github.com/cyiafn/flight_information_system/server/utils/bytes"
	"github.com/pkg/errors"
)

/*
Note: no map support, no nested array support
*/

// constants for primitives
const (
	intSize     = 8
	int64Size   = 8
	uint8Size   = 1
	int32Size   = 4
	float64Size = 8
)

var (
	// stringTerminator (\0)
	stringTerminator = []byte("\000")[0]
)

// Unmarshal a byte array to a structure based on the structure outlined in the report.
// Note, nested arrays are not implemented, maps are not implemented, some primitives unused are not implemented as well
// This uses a lot of runtime evaluation with some meta programming, it is not as performant as the standard marshalling library.
// we keep a ptr while unmarshalling to indicate the index of the byte we are on
func Unmarshal(request []byte, v any) error {
	var err error

	// gets the actual type of the structure
	reflectValue := reflect.ValueOf(v)
	reflectElem := reflectValue.Elem()

	ptr := 0

	// if v != structure, its an error
	if reflectElem.Kind() != reflect.Struct {
		logs.Error("value passed in is not of structure type")
		return custom_errors.NewMarshallerError(errors.Errorf("value passed in is not of structure type"))
	}

	// for each field, we populate the data in sequence.
	for i := 0; i < reflectElem.NumField(); i++ {
		field := reflectElem.FieldByName(reflectElem.Type().Field(i).Name)
		// if we are able to manipulate the field
		if field.IsValid() && field.CanSet() {
			fieldKind := reflectElem.Type().Field(i).Type.Kind()
			// if it is an interface, we need to evaluate further whats the actual type
			if fieldKind == reflect.Interface {
				if field.IsNil() {
					continue
				}
				fieldKind = reflect.TypeOf(field.Interface()).Elem().Kind()
			}
			// for each type we unmarshal
			switch fieldKind {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Float64, reflect.String:
				ptr, err = unmarshalPrimitive(request, fieldKind, field, ptr)
			case reflect.Slice:
				ptr, err = unmarshalArray(request, field, reflectElem.Type().Field(i).Type.Elem().Kind(), ptr)
			case reflect.Struct:
				ptr, err = unmarshalStruct(request, field, ptr)
			default:
				logs.Error("unimplemented type: %v", fieldKind)
				return custom_errors.NewMarshallerError(errors.Errorf("unimplemented type, type: %v", fieldKind))
			}
		}
	}
	return err
}

// unmarshalArray unmarshals part of the byte array to an array
func unmarshalArray(request []byte, field reflect.Value, elementType reflect.Kind, ptr int) (int, error) {
	// gets the length of the array
	sizeOfSlice := int(bytes.ToInt64(request[ptr : ptr+int64Size]))
	ptr += int64Size

	// for each different element type, we convert the value from []byte to the actual data type
	// recursively unmarshals for structure type
	switch elementType {
	case reflect.Int:
		slice := reflect.MakeSlice(reflect.TypeOf([]int{}), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			slice.Index(i).Set(reflect.ValueOf(int(bytes.ToInt64(request[ptr : ptr+intSize]))))
			ptr += intSize
		}
		field.Set(slice)
	case reflect.Int32:
		slice := reflect.MakeSlice(reflect.TypeOf([]int32{}), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			slice.Index(i).Set(reflect.ValueOf(bytes.ToInt32(request[ptr : ptr+int32Size])))
			ptr += int32Size
		}
		field.Set(slice)

	case reflect.Int64:
		slice := reflect.MakeSlice(reflect.TypeOf([]int64{}), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			slice.Index(i).Set(reflect.ValueOf(bytes.ToInt64(request[ptr : ptr+int64Size])))
			ptr += int64Size
		}
		field.Set(slice)
	case reflect.Uint8:
		slice := reflect.MakeSlice(reflect.TypeOf([]uint8{}), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			slice.Index(i).Set(reflect.ValueOf(request[ptr]))
			ptr += uint8Size
		}
		field.Set(slice)
	case reflect.Float64:
		slice := reflect.MakeSlice(reflect.TypeOf([]float64{}), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			slice.Index(i).Set(reflect.ValueOf(bytes.ToFloat64(request[ptr : ptr+float64Size])))
			ptr += float64Size
		}
		field.Set(slice)
	case reflect.String:
		slice := reflect.MakeSlice(reflect.TypeOf([]string{}), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			start := ptr
			for ; start < len(request); start++ {
				if request[start] == stringTerminator { // we get all bytes until the string terminator \0
					break
				}
			}

			if start == ptr {
				logs.Fatal("buffer for UDP datagram not big enough for a single request")
			}

			slice.Index(i).Set(reflect.ValueOf(string(request[ptr:start])))
			ptr = start + 1
		}
		field.Set(slice)
	case reflect.Struct:
		slice := reflect.MakeSlice(reflect.SliceOf(field.Type().Elem()), sizeOfSlice, sizeOfSlice)
		for i := 0; i < sizeOfSlice; i++ {
			ind := slice.Index(i)
			var err error
			ptr, err = unmarshalStruct(request, ind, ptr) // recursively unmarshals the structure for each index
			if err != nil {
				return 0, err
			}
		}
		field.Set(slice)
	default:
		logs.Error("unimplemented type: %v", elementType)
		return 0, custom_errors.NewMarshallerError(errors.Errorf("unimplemented type, type: %v", elementType))
	}
	return ptr, nil
}

// unmarshalStruct is mostly similar to unmarshal
func unmarshalStruct(request []byte, reflectValue reflect.Value, ptr int) (int, error) {
	var err error

	// for each field we unmarshal based on the type
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.FieldByName(reflectValue.Type().Field(i).Name)
		if field.IsValid() && field.CanSet() {
			fieldKind := reflectValue.Type().Field(i).Type.Kind()
			switch fieldKind {
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint8, reflect.Float64, reflect.String:
				ptr, err = unmarshalPrimitive(request, fieldKind, field, ptr)
				if err != nil {
					return 0, err
				}
			case reflect.Slice:
				ptr, err = unmarshalArray(request, field, reflectValue.Type().Field(i).Type.Elem().Kind(), ptr)
				if err != nil {
					return 0, err
				}
			case reflect.Struct:
				ptr, err = unmarshalStruct(request, field, ptr)
				if err != nil {
					return 0, err
				}
			default:
				logs.Error("unimplemented type: %v", fieldKind)
				return 0, custom_errors.NewMarshallerError(errors.Errorf("unimplemented type, type: %v", fieldKind))
			}
		}
	}
	return ptr, nil
}

// unmarshalPrimitive unmarshals  each primitive into v
func unmarshalPrimitive(request []byte, fieldKind reflect.Kind, field reflect.Value, ptr int) (int, error) {
	switch fieldKind {
	case reflect.Int, reflect.Int64:
		field.SetInt(bytes.ToInt64(request[ptr : ptr+intSize]))
		ptr += intSize
	case reflect.Int32:
		field.SetInt(int64(bytes.ToInt32(request[ptr : ptr+int32Size])))
		ptr += int32Size
	case reflect.Uint8:
		field.SetUint(uint64(request[ptr]))
		ptr += uint8Size
	case reflect.Float64:
		field.SetFloat(bytes.ToFloat64(request[ptr : ptr+float64Size]))
		ptr += float64Size
	case reflect.String:
		endString := ptr
		for ; endString < len(request); endString++ {
			if request[endString] == stringTerminator {
				break
			}
		}

		field.SetString(string(request[ptr:endString]))
		ptr = endString + 1
	default:
		logs.Error("unimplemented type: %v", fieldKind)
		return 0, custom_errors.NewMarshallerError(errors.Errorf("unimplemented type, type: %v", fieldKind))
	}
	return ptr, nil
}
