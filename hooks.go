package mapstructure_hooks

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func StringToNumericHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() == reflect.Int {
			i, err := strconv.ParseInt(data.(string), 10, 64)
			return int(i), err
		}

		if t.Kind() == reflect.Int8 {
			i, err := strconv.ParseInt(data.(string), 10, 8)
			return int8(i), err
		}

		if t.Kind() == reflect.Int16 {
			i, err := strconv.ParseInt(data.(string), 10, 16)
			return int16(i), err
		}

		if t.Kind() == reflect.Int32 {
			i, err := strconv.ParseInt(data.(string), 10, 32)
			return int32(i), err
		}

		if t.Kind() == reflect.Int64 {
			i, err := strconv.ParseInt(data.(string), 10, 64)
			return int64(i), err
		}

		if t.Kind() == reflect.Uint {
			i, err := strconv.ParseUint(data.(string), 10, 64)
			return uint(i), err
		}

		if t.Kind() == reflect.Uint8 {
			i, err := strconv.ParseUint(data.(string), 10, 8)
			return uint8(i), err
		}

		if t.Kind() == reflect.Uint16 {
			i, err := strconv.ParseUint(data.(string), 10, 16)
			return uint16(i), err
		}

		if t.Kind() == reflect.Uint32 {
			i, err := strconv.ParseUint(data.(string), 10, 32)
			return uint32(i), err
		}

		if t.Kind() == reflect.Uint64 {
			i, err := strconv.ParseUint(data.(string), 10, 64)
			return uint64(i), err
		}

		if t.Kind() == reflect.Float32 {
			f, err := strconv.ParseFloat(data.(string), 32)
			return float32(f), err
		}

		if t.Kind() == reflect.Float64 {
			f, err := strconv.ParseFloat(data.(string), 64)
			return float64(f), err
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
		StringToNumericHookFunc(),
		StringToStringSliceHookFunc(),
		StringToURLHookFunc(),
		OptionHookFunc(),
	}
}
