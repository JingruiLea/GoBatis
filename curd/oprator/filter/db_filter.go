package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type DBFilter interface {
	Filter(db *gorm.DB) (*gorm.DB, error)
	SetValue(v reflect.Value)
	NeedValue() bool
	Check(in []reflect.Type, i int) (string, int)
}

type DefaultFilter struct {
	Active bool
}

func (d DefaultFilter) NeedValue() bool {
	return d.Active
}
