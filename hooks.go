package mapstructure_hooks

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	protoreflectEnumType = reflect.TypeOf((*protoreflect.Enum)(nil)).Elem()
)

func StringToProtoEnumHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if !t.Implements(protoreflectEnumType) {
			return data, nil
		}

		var (
			name     = protoreflect.Name(data.(string))
			desc     = reflect.New(t).Elem().MethodByName("Descriptor").Call(nil)[0]
			enumName = desc.MethodByName("Name").Call(nil)[0]
			values   = desc.MethodByName("Values").Call(nil)[0]
			value    = values.MethodByName("ByName").Call([]reflect.Value{reflect.ValueOf(name)})[0]
		)

		if value.IsNil() {
			return nil, fmt.Errorf("invalid %s value: %s", enumName.Interface(), name)
		}

		var num = value.MethodByName("Number").Call(nil)[0]
		return num.Interface(), nil
	}
}

func StringToNumericHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		var (
			s          = data.(string)
			multiplier = 1
		)

		switch {
		case strings.HasSuffix(s, "Ki"):
			multiplier = 1024
			s = strings.TrimSuffix(s, "Ki")
		case strings.HasSuffix(s, "Mi"):
			multiplier = 1024 * 1024
			s = strings.TrimSuffix(s, "Mi")
		case strings.HasSuffix(s, "Gi"):
			multiplier = 1024 * 1024 * 1024
			s = strings.TrimSuffix(s, "Gi")
		case strings.HasSuffix(s, "K"):
			multiplier = 1000
			s = strings.TrimSuffix(s, "K")
		case strings.HasSuffix(s, "M"):
			multiplier = 1000 * 1000
			s = strings.TrimSuffix(s, "M")
		case strings.HasSuffix(s, "G"):
			multiplier = 1000 * 1000 * 1000
			s = strings.TrimSuffix(s, "G")
		}

		if t.Kind() == reflect.Int {
			i, err := strconv.ParseInt(s, 10, 64)
			return int(i) * int(multiplier), err
		}

		if t.Kind() == reflect.Int8 {
			i, err := strconv.ParseInt(s, 10, 8)
			return int8(i) * int8(multiplier), err
		}

		if t.Kind() == reflect.Int16 {
			i, err := strconv.ParseInt(s, 10, 16)
			return int16(i) * int16(multiplier), err
		}

		if t.Kind() == reflect.Int32 {
			i, err := strconv.ParseInt(s, 10, 32)
			return int32(i) * int32(multiplier), err
		}

		if t.Kind() == reflect.Int64 {
			i, err := strconv.ParseInt(s, 10, 64)
			return int64(i) * int64(multiplier), err
		}

		if t.Kind() == reflect.Uint {
			i, err := strconv.ParseUint(s, 10, 64)
			return uint(i) * uint(multiplier), err
		}

		if t.Kind() == reflect.Uint8 {
			i, err := strconv.ParseUint(s, 10, 8)
			return uint8(i) * uint8(multiplier), err
		}

		if t.Kind() == reflect.Uint16 {
			i, err := strconv.ParseUint(s, 10, 16)
			return uint16(i) * uint16(multiplier), err
		}

		if t.Kind() == reflect.Uint32 {
			i, err := strconv.ParseUint(s, 10, 32)
			return uint32(i) * uint32(multiplier), err
		}

		if t.Kind() == reflect.Uint64 {
			i, err := strconv.ParseUint(s, 10, 64)
			return uint64(i) * uint64(multiplier), err
		}

		if t.Kind() == reflect.Float32 {
			f, err := strconv.ParseFloat(s, 32)
			return float32(f) * float32(multiplier), err
		}

		if t.Kind() == reflect.Float64 {
			f, err := strconv.ParseFloat(s, 64)
			return float64(f) * float64(multiplier), err
		}

		return data, nil
	}
}

func StringToBoolHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if (f.Kind() != reflect.String) || (t.Kind() != reflect.Bool) {
			return data, nil
		}

		return strconv.ParseBool(data.(string))
	}
}

func StringToStringSliceHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if (f.Kind() != reflect.String) || (t != reflect.TypeOf([]string{})) {
			return data, nil
		}

		collection := make([]string, 0)

		for _, item := range strings.Split(data.(string), ",") {
			collection = append(collection, item)
		}

		return collection, nil
	}
}

func StringToURLHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if (f.Kind() != reflect.String) || (t != reflect.TypeOf(url.URL{})) {
			return data, nil
		}

		url, err := url.Parse(data.(string))

		if err != nil {
			return nil, err
		}

		return *url, nil
	}
}

func OptionHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if !strings.HasPrefix(t.String(), "mo.Option[") {
			return data, nil
		}

		b, err := json.Marshal(data)

		if err != nil {
			return nil, err
		}

		var inst = reflect.New(t).Interface()

		if err := json.Unmarshal(b, inst); err != nil {
			return nil, err
		}

		return inst, nil
	}
}

func All() []mapstructure.DecodeHookFunc {
	return []mapstructure.DecodeHookFunc{
		StringToBoolHookFunc(),
		StringToProtoEnumHookFunc(),
		StringToNumericHookFunc(),
		StringToStringSliceHookFunc(),
		StringToURLHookFunc(),
		OptionHookFunc(),
	}
}
