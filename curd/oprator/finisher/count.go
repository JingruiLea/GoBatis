package finisher

import (
	"reflect"

	"gorm.io/gorm"
)

type Count struct {
}

func (c *Count) Finish(db *gorm.DB, dbFunc *DBFunc) (res []reflect.Value, err error) {
	var count int64
	err = db.Count(&count).Error
	res = append(res, reflect.ValueOf(count))
	return res, err
}

func (c *Count) Fill([]reflect.Value) {

}

func (c *Count) CheckRes(o []reflect.Type) string {
	if len(o) == 2 && o[1].Kind() == reflect.Int64 {
		return ""
	}
	return "Count方法应该返回(int64, error)"
}
