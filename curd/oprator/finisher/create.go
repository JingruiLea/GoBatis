package finisher

import (
	"reflect"

	"gorm.io/gorm"
)

type Create struct {
	Model interface{}
}

func (c *Create) Finish(db *gorm.DB, dbFunc *DBFunc) (res []reflect.Value, err error) {
	if dbFunc.HasBO && dbFunc.IsBO(c.Model) {
		c.Model = dbFunc.BO2Model(c.Model)
	}
	err = db.Create(c.Model).Error
	return res, err
}

func (c *Create) Fill(vs []reflect.Value) {
	c.Model = vs[0].Interface()
}

func (c *Create) CheckRes(o []reflect.Type) string {
	if len(o) == 1 {
		return ""
	}
	return "Create方法应该返回(error)"
}
