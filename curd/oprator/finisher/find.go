package finisher

import (
	"reflect"

	"gorm.io/gorm"
)

func (f *Find) Finish(db *gorm.DB, dbFunc *DBFunc) (res []reflect.Value, err error) {
	var out reflect.Value
	if f.First {
		out = reflect.New(dbFunc.ModelT)
		err = db.First(out.Interface()).Error
		if f.ResBO {
			bo := dbFunc.Model2BO(out.Interface())
			res = append(res, reflect.ValueOf(bo))
		} else {
			res = append(res, out)
		}
	} else {
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.New(dbFunc.ModelT).Type()), 0, 0)
		out := reflect.New(slice.Type())
		out.Elem().Set(slice)
		err = db.Find(out.Interface()).Error
		if f.ResBO {
			boArr := dbFunc.Model2BO(out.Elem().Interface())
			res = append(res, reflect.ValueOf(boArr))
		} else {
			res = append(res, out.Elem())
		}
	}
	return res, err
}

type Find struct {
	First bool
	ResBO bool
}

func (f *Find) Fill(vs []reflect.Value) {

}

func (f *Find) CheckRes(o []reflect.Type) string {
	if f.First {
		if len(o) == 2 && o[0].Kind() == reflect.Ptr && o[0].Elem().Kind() == reflect.Struct {
			return ""
		}
		return "Get方法应该返回(*model, error)"
	} else {
		if len(o) == 2 && o[0].Kind() == reflect.Slice && o[0].Elem().Kind() == reflect.Ptr && o[0].Elem().Elem().Kind() == reflect.Struct {
			return ""
		}
		return "MGet方法应该返回([]*model, error)"
	}
}
