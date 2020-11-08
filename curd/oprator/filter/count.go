package filter

import (
	"reflect"

	"gorm.io/gorm"
)

type Count struct {
	V     reflect.Value
	Count int64
	AsRes bool
	DefaultFilter
}

func (c *Count) Filter(db *gorm.DB) (res *gorm.DB, err error) {
	if c.AsRes {
		err = db.Count(&c.Count).Error
		c.V = reflect.ValueOf(&c.Count)
		return db, err
	}
	if c.V.IsValid() && !c.V.IsNil() {
		err = db.Count(&c.Count).Error
		c.V.Elem().Set(reflect.ValueOf(c.Count))
	}
	return db, err
}

func (c *Count) SetValue(v reflect.Value) {
	c.V = v
}

func (c *Count) Check(in []reflect.Type, i int) (string, int) {
	if c.AsRes {
		return "", i
	}
	t := in[i]
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Int64 {
		return "", i + 1
	}
	return "Count对应的参数必须为 *int64", i + 1
}
