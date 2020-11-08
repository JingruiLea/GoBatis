package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Limit struct {
	V int
	DefaultFilter
}

func (l *Limit) Filter(db *gorm.DB) (*gorm.DB, error) {
	if l.V == 0 {
		return db, nil
	}
	return db.Limit(l.V), nil
}

func (l *Limit) SetValue(v reflect.Value) {
	l.V = int(v.Int())
}

func (l *Limit) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() == reflect.Int ||
		t.Kind() == reflect.Int32 ||
		t.Kind() == reflect.Int64 {
		return "", i + 1
	}
	return "Limit对应的参数必须为 int", i + 1
}
