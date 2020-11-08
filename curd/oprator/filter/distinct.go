package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Distinct struct {
	V string
	DefaultFilter
}

func (l *Distinct) Filter(db *gorm.DB) (*gorm.DB, error) {
	if l.V == "" {
		return db, nil
	}
	return db.Distinct(l.V), nil
}

func (l *Distinct) SetValue(v reflect.Value) {
	l.V = v.String()
}

func (l *Distinct) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() == reflect.String {
		return "", i + 1
	}
	return "Distinct对应的参数必须为 string", i + 1
}
