package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Order struct {
	V string
	DefaultFilter
}

func (o *Order) Filter(db *gorm.DB) (*gorm.DB, error) {
	if o.V == "" {
		return db, nil
	}
	return db.Order(o.V), nil
}

func (o *Order) SetValue(v reflect.Value) {
	switch v.Interface().(type) {
	case string:
		o.V = v.String()
	}
}

func (o *Order) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() == reflect.String {
		return "", i + 1
	}
	return "Order对应的参数必须为 string", i + 1
}
