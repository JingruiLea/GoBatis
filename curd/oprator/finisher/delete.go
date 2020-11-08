package finisher

import (
	"reflect"

	"gorm.io/gorm"
)

type Delete struct {
}

func (d *Delete) Finish(db *gorm.DB, dbFunc *DBFunc) (res []reflect.Value, err error) {
	err = db.Delete(reflect.New(dbFunc.ModelT).Interface()).Error
	if dbFunc.T.NumOut() == 2 {
		res = append(res, reflect.ValueOf(db.RowsAffected))
	}
	return res, err
}

func (d *Delete) Fill(value []reflect.Value) {}

func (d *Delete) CheckRes(o []reflect.Type) string {
	if len(o) == 1 {
		return ""
	}
	if len(o) == 2 && o[1].Kind() == reflect.Int64 {
		return ""
	}
	return "Delete方法应该返回(error)或者(rowsAffected int64, err error)"
}
