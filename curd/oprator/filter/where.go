package filter

import (
	"reflect"
	"strings"

	"github.com/JingruiLea/gobatis/curd/utils"

	"gorm.io/gorm"
)

type Where struct {
	Col       string
	Op        string
	Obj       interface{}
	Or        bool
	Essential bool
	Model     interface{}
	DefaultFilter
}

func (sel *Where) SetValue(value reflect.Value) {
	if sel.Col != "" {
		if (value.Kind() == reflect.Array ||
			value.Kind() == reflect.Slice) && sel.Op == "=" {
			sel.Op = "in"
		}
		sel.Obj = value.Interface()
	} else {
		sel.Model = value.Interface()
	}
}

func (sel *Where) Filter(db *gorm.DB) (*gorm.DB, error) {
	var res *gorm.DB
	if sel.Model != nil {
		res = db.Where(sel.Model)
		return res, nil
	}
	if !sel.HasSetObj() {
		return db, nil
	}
	sel.Col = strings.TrimSpace(sel.Col)
	if sel.Or {
		res = db.Or(sel.Col+" "+sel.Op+" ?", sel.Obj)
	} else {
		res = db.Where(sel.Col+" "+sel.Op+" ?", sel.Obj)
	}
	return res, nil
}

func (sel *Where) HasSetObj() bool {
	s, ok := sel.Obj.(string)
	if !ok {
		return true
	}
	if s == "?" {
		return false
	}
	return true
}

func (sel *Where) HasValue() bool {
	if sel.Model != nil {
		return true
	}
	if !sel.HasSetObj() {
		return false
	}
	if !utils.IsParamsValid(sel.Obj) {
		return false
	}
	return true
}

func (sel *Where) Check(in []reflect.Type, i int) (string, int) {
	return "", i + 1
}
