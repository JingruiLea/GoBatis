package finisher

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/JingruiLea/gobatis/curd/constants"

	"github.com/JingruiLea/gobatis/curd/oprator/filter"

	"gorm.io/gorm"
)

var GetWriteDB func(context.Context) *gorm.DB
var GetReadDB func(context.Context) *gorm.DB

type DBFunc struct {
	ModelT     reflect.Type
	BOT        reflect.Type
	HasBO      bool
	Name       string
	StructT    reflect.Type
	StructV    reflect.Value
	T          reflect.Type
	V          reflect.Value
	Tag        reflect.StructTag
	StructName string
	FuncType   constants.FuncType

	Filters  []filter.DBFilter
	Finisher DBFinisher
}

func (f *DBFunc) IsBO(o interface{}) bool {
	ot := reflect.ValueOf(o).Type()
	if ot.Kind() == reflect.Array || ot.Kind() == reflect.Slice {
		ot = ot.Elem().Elem()
	} else {
		ot = ot.Elem()
	}
	return ot.ConvertibleTo(f.BOT)
}

func (f *DBFunc) IsBOT(ot reflect.Type) bool {
	if !f.HasBO {
		return false
	}
	if ot.Kind() == reflect.Array || ot.Kind() == reflect.Slice {
		ot = ot.Elem().Elem()
	}
	return ot.ConvertibleTo(f.BOT)
}

func (f *DBFunc) Model2BO(o interface{}) interface{} {
	ov := reflect.ValueOf(o)
	if ov.Kind() == reflect.Array || ov.Kind() == reflect.Slice {
		res := reflect.MakeSlice(reflect.SliceOf(reflect.New(f.BOT).Type()), ov.Len(), ov.Len())
		for i := 0; i < ov.Len(); i++ {
			model := ov.Index(i)
			method := model.MethodByName("ToBO")
			refRes := method.Call([]reflect.Value{})
			res.Index(i).Set(refRes[0])
		}
		return res.Interface()
	} else {
		method := ov.MethodByName("ToBO")
		refRes := method.Call([]reflect.Value{})
		return refRes[0].Interface()
	}
}

func (f *DBFunc) BO2Model(o interface{}) interface{} {
	ov := reflect.ValueOf(o)
	if ov.Kind() == reflect.Array || ov.Kind() == reflect.Slice {
		res := reflect.MakeSlice(reflect.SliceOf(reflect.New(f.ModelT).Type()), ov.Len(), ov.Len())
		for i := 0; i < ov.Len(); i++ {
			model := reflect.New(f.ModelT)
			method := model.MethodByName("FromBO")
			method.Call([]reflect.Value{ov.Index(i)})
			res.Index(i).Set(model)
		}
		return res.Interface()
	} else {
		model := reflect.New(f.ModelT)
		method := model.MethodByName("FromBO")
		method.Call([]reflect.Value{ov})
		return model.Interface()
	}
}

func (f *DBFunc) Check() (err error) {
	numOut := f.T.NumOut()
	out := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		out[i] = f.T.Out(i)
	}
	numIn := f.T.NumIn()
	in := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		in[i] = f.T.In(i)
	}

	errMsg := ""
	defer func() {
		if errMsg != "" {
			err = errors.New(fmt.Sprintf("%s.%s fill error.%s", f.StructName, f.Name, errMsg))
		}
	}()
	//基本参数检查
	ctxt := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !in[0].Implements(ctxt) {
		errMsg = fmt.Sprintf("第一个参数必须为context")
		return
	}
	if !in[1].ConvertibleTo(reflect.TypeOf(&gorm.DB{})) {
		errMsg = fmt.Sprintf("第二个参数必须为*gorm.db")
		return
	}
	errt := reflect.TypeOf((*error)(nil)).Elem()
	if !out[numOut-1].Implements(errt) {
		errMsg = fmt.Sprintf("最后一个返回值必须为error")
		return
	}
	modelT := reflect.New(f.ModelT).Type()
	_, hasFrom := modelT.MethodByName("FromBO")
	_, hasTo := modelT.MethodByName("ToBO")
	if f.HasBO && !(hasFrom && hasTo) {
		errMsg = fmt.Sprintf("使用BO的func必须定义FromBO和ToBO方法")
		return
	}

	var resCount *filter.Count
	//参数数量检查
	argIndex := 2
	for _, dbFilter := range f.Filters {
		if c, ok := dbFilter.(*filter.Count); ok && c.AsRes {
			resCount = c
		}
		if dbFilter.NeedValue() {
			errMsg, _ = dbFilter.Check(in, argIndex)
			if errMsg != "" {
				return
			}
			argIndex++
		}
	}
	u, ok := f.Finisher.(*Update)
	if ok {
		argIndex += u.Need()
	}
	_, ok = f.Finisher.(*Create)
	if ok {
		argIndex++
	}
	if argIndex < numIn {
		errMsg = fmt.Sprintf("参数数量不匹配, 需要%d, 提供了%d", argIndex, numIn)
		return
	}

	if resCount != nil {
		if out[len(out)-2].Kind() != reflect.Int64 {
			errMsg = fmt.Sprintf("count:\"res\" 需要倒数第二个返回值为int64类型")
			return
		}
		out = append(out[0:len(out)-2], errt)
	}
	//返回类型检查
	errMsg = f.Finisher.CheckRes(out)
	return err
}

func (f *DBFunc) Exec(args []reflect.Value) (res []reflect.Value, err error) {
	ctx := args[0].Interface().(context.Context)
	defer func() {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = f.GetRecordNotFound(ctx, err)
		}
	}()
	db := args[1].Interface().(*gorm.DB)
	if db == nil {
		switch f.FuncType {
		case constants.Create, constants.Delete, constants.Update, constants.UpdateBy, constants.DeleteBy:
			db = GetWriteDB(ctx)
		default:
			db = GetReadDB(ctx)
		}
	}
	argIndex := 2
	var page *filter.Page
	var resCount *filter.Count
	for _, dbFilter := range f.Filters {
		if dbFilter.NeedValue() {
			dbFilter.SetValue(args[argIndex])
			argIndex++
		}
	}
	for _, dbFilter := range f.Filters {
		if s, ok := dbFilter.(*filter.Where); ok {
			if !s.Essential && !s.HasValue() {
				continue
			}
		}
		if p, ok := dbFilter.(*filter.Page); ok {
			page = p
			continue
		}
		if c, ok := dbFilter.(*filter.Count); ok && c.AsRes {
			resCount = c
		}
		db, err = dbFilter.Filter(db)
		if err != nil {
			return res, err
		}
	}
	if page != nil { //page放在最后
		db, err = page.Filter(db)
		if err != nil {
			return res, err
		}
	}

	f.Finisher.Fill(args[argIndex:])
	res, err = f.Finisher.Finish(db, f)
	if err != nil {
		return res, err
	}
	if resCount != nil {
		res = append(res, resCount.V)
	}
	return res, err
}

type RecordNotFoundImp interface {
	RecordNotFound(ctx context.Context, structName string, funcName string) error
}

func (f *DBFunc) GetRecordNotFound(ctx context.Context, err error) error {
	recordNotFoundT := reflect.TypeOf((*RecordNotFoundImp)(nil)).Elem()
	if f.StructT.Implements(recordNotFoundT) {
		return f.StructV.Interface().(RecordNotFoundImp).RecordNotFound(ctx, f.StructName, f.Name)
	}
	return err
}

type BOConvertible interface {
	FromBO(BO interface{})
	ToBO() interface{}
}
