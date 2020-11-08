package filter

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Lock struct {
	V bool
	DefaultFilter
}

func (l *Lock) Filter(db *gorm.DB) (*gorm.DB, error) {
	if l.V {
		return db.Clauses(clause.Locking{Strength: "UPDATE"}), nil
	}
	return db, nil
}

func (l *Lock) SetValue(v reflect.Value) {
	l.V = v.Bool()
}

func (l *Lock) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() == reflect.Bool {
		return "", i + 1
	}
	return "Group对应的参数必须为 string", i + 1
}
