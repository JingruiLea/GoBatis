package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Group struct {
	V string
	DefaultFilter
}

func (l *Group) Filter(db *gorm.DB) (*gorm.DB, error) {
	if l.V == "" {
		return db, nil
	}
	return db.Group(l.V), nil
}

func (l *Group) SetValue(v reflect.Value) {
	l.V = v.String()
}

func (l *Group) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() == reflect.String {
		return "", i + 1
	}
	return "Group对应的参数必须为 string", i + 1
}
