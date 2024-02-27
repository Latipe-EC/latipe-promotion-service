package mapper

import (
	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"
)

// BindingStruct - biding struct to struct
func BindingStructGrpc(src interface{}, desc proto.Message) error {
	byteSrc, err := sonic.Marshal(src)
	if err != nil {
		return err
	}
	// binding to desc
	if err := protojson.Unmarshal(byteSrc, desc); err != nil {
		return err
	}
	return nil
}

// BindingStruct - biding struct to struct
func BindingStruct(src interface{}, desc interface{}) error {
	// convert to byte
	byteSrc, err := sonic.Marshal(src)
	if err != nil {
		return err
	}
	// binding to desc
	err = sonic.Unmarshal(byteSrc, &desc)
	if err != nil {
		return err
	}
	return nil
}

func BindingAndValidate[T any](detail interface{}, validator func(interface{}) error) (T, error) {
	var model T
	if err := BindingStruct(detail, &model); err != nil {
		return model, err
	}

	if err := validator(model); err != nil {
		return model, err
	}
	return model, nil
}

func StructToMap(input interface{}, ignoreNilFiled bool) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(input)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get("json") == "" {
			continue
		}

		fv := v.Field(i)
		if ignoreNilFiled && fv.Type().Kind() == reflect.Pointer && fv.IsNil() {
			continue
		}

		if fv.Kind() == reflect.Pointer {
			fv = fv.Elem()
		}

		value := fv.Interface()
		result[field.Tag.Get("json")] = value
	}
	return result
}
