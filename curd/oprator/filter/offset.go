package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Offset struct {
	V int
	DefaultFilter
}

func (l *Offset) Filter(db *gorm.DB) (*gorm.DB, error) {
	if l.V == 0 {
		return db, nil
	}
	return db.Offset(l.V), nil
}

func (l *Offset) SetValue(v reflect.Value) {
	l.V = int(v.Int())
}

func (l *Offset) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() == reflect.Int ||
		t.Kind() == reflect.Int32 ||
		t.Kind() == reflect.Int64 {
		return "", i + 1
	}
	return "Offset对应的参数必须为 int", i + 1
}
