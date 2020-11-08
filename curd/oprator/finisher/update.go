package finisher

import (
	"reflect"
	"strings"

	"github.com/JingruiLea/gobatis/curd/utils"

	"gorm.io/gorm"
)

type Update struct {
	Updates     map[string]interface{}
	ks          []string
	vs          []interface{}
	Obj         interface{}
	UpdateModel bool
	NeedC       int
}

type Expr struct {
	V     string
	Count int
}

func (u *Update) Finish(db *gorm.DB, dbFunc *DBFunc) (res []reflect.Value, err error) {
	if u.Obj != nil {
		if dbFunc.HasBO && dbFunc.IsBO(u.Obj) {
			u.Obj = dbFunc.BO2Model(u.Obj)
		}
		err = db.Updates(u.Obj).Error
	} else {
		if len(u.Updates) != 0 {
			err = db.Updates(u.Updates).Error
		}
	}
	if dbFunc.T.NumOut() == 2 {
		res = append(res, reflect.ValueOf(db.RowsAffected))
	}
	return res, err
}

func (u *Update) CheckRes(o []reflect.Type) string {
	if len(o) == 1 {
		return ""
	}
	if len(o) == 2 && o[0].Kind() == reflect.Int64 {
		return ""
	}
	return "update方法应该返回(error)或者(rowsAffected int64, err error)"
}

func (u *Update) Need() int {
	if u.UpdateModel {
		return 1
	}
	count := 0
	for _, v := range u.vs {
		if v == nil {
			count++
		}
		if e, ok := v.(Expr); ok {
			count += e.Count
		}
	}
	return count
}

func (u *Update) Fill(vs []reflect.Value) {
	if u.UpdateModel {
		u.Obj = vs[0].Interface()
		return
	}
	u.Updates = make(map[string]interface{})
	index := 0
	for i, v := range u.vs {
		k := u.ks[i]
		if v == nil {
			v = vs[index].Interface()
			index++
		}
		if e, ok := v.(Expr); ok {
			params := make([]interface{}, 0)
			for i := 0; i < e.Count; i++ {
				params = append(params, vs[index].Interface())
				index++
			}
			v = gorm.Expr(e.V, params...)
		}
		if strings.HasSuffix(k, "?") {
			k = strings.TrimRight(k, "?")
			if !utils.IsParamsValid(v) {
				continue
			}
		}
		u.Updates[k] = v
	}
}

func (u *Update) Set(col string, v interface{}) {
	u.ks = append(u.ks, col)
	u.vs = append(u.vs, v)
}

func (u *Update) HasSetCol(col string) bool {
	s, ok := u.Updates[col]
	if !ok {
		return false
	}
	str, ok := s.(string)
	if !ok {
		return true
	}
	if str == "?" {
		return false
	}
	return true
}
