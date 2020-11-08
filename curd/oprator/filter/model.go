package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Model struct {
	Table string
	Model reflect.Type
	DefaultFilter
}

func (m *Model) Filter(db *gorm.DB) (*gorm.DB, error) {
	if m.Model != nil {
		orm := reflect.New(m.Model).Interface()
		return db.Model(orm), nil
	}
	return db.Table(m.Table), nil
}

func (m *Model) SetValue(v reflect.Value) {
	switch v.Interface().(type) {
	case string:
		m.Table = v.String()
	default:
		m.Model = v.Type()
	}
}
func (m *Model) Check(in []reflect.Type, i int) (string, int) {
	return "", i
}
