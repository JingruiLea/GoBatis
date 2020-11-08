package curd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/segmentio/encoding/json"

	"github.com/JingruiLea/gobatis/curd/parser"

	"github.com/JingruiLea/gobatis/curd/utils"

	"github.com/JingruiLea/gobatis/curd/oprator/filter"
	"github.com/JingruiLea/gobatis/curd/oprator/finisher"

	"github.com/JingruiLea/gobatis/curd/constants"

	"gorm.io/gorm"
)

type Create func(context.Context, *gorm.DB, interface{}) error
type Delete func(ctx context.Context, db *gorm.DB, v interface{}) error

func FillForDebug(dalIns interface{}, funcName string) error {
	dalV := reflect.ValueOf(dalIns)
	if dalV.Kind() != reflect.Ptr || dalV.IsNil() || dalV.IsZero() {
		return errors.New("类型错误, 传入的dal Object 必须是非空指针")
	}
	dalT := reflect.TypeOf(dalIns).Elem()
	dalV = dalV.Elem()
	modelT, ok := dalT.FieldByName("Model")
	if !ok {
		modelT, ok = dalT.FieldByName("model")
		if !ok {
			return errors.New("dal结构体未定义Model域")
		}
	}
	for i := 0; i < dalT.NumField(); i++ {
		f := dalT.Field(i)
		fV := dalV.Field(i)
		if f.Name == funcName {
			dbFunc := &finisher.DBFunc{
				Name:       f.Name,
				ModelT:     modelT.Type.Elem(),
				StructT:    dalT,
				StructV:    dalV,
				T:          f.Type,
				V:          fV,
				Tag:        f.Tag,
				StructName: dalT.Name(),
			}
			err := FillFunc(dbFunc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AutoFill(dalIns ...interface{}) error {
	for _, in := range dalIns {
		err := autoFillOne(in)
		if err != nil {
			return err
		}
	}
	return nil
}

func autoFillOne(dalIns interface{}) error {
	dalV := reflect.ValueOf(dalIns)
	if dalV.Kind() != reflect.Ptr || dalV.IsNil() || dalV.IsZero() {
		return errors.New("类型错误, 传入的dal Object 必须是非空指针")
	}
	dalT := reflect.TypeOf(dalIns).Elem()
	dalV = dalV.Elem()
	modelT, ok := dalT.FieldByName("Model")
	if !ok {
		modelT, ok = dalT.FieldByName("model")
		if !ok {
			return errors.New("dal结构体未定义Model域")
		}
	}
	BOT, hasBO := dalT.FieldByName("BO")
	if !hasBO {
		BOT, hasBO = dalT.FieldByName("bo")
	}
	for i := 0; i < dalT.NumField(); i++ {
		f := dalT.Field(i)
		fV := dalV.Field(i)
		if f.Name == "Model" || f.Name == "BO" || f.Name == "model" || f.Name == "bo" {
			continue
		}
		dbFunc := &finisher.DBFunc{
			ModelT:     modelT.Type.Elem(),
			HasBO:      hasBO,
			Name:       f.Name,
			StructT:    dalT,
			StructV:    dalV,
			T:          f.Type,
			V:          fV,
			Tag:        f.Tag,
			StructName: dalT.Name(),
			FuncType:   "",
			Filters:    nil,
			Finisher:   nil,
		}
		if hasBO {
			dbFunc.BOT = BOT.Type.Elem()
		}
		log.Printf("fill %s.%s\n", dalT.Name(), f.Name)
		err := FillFunc(dbFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckFuncLength(ft constants.FuncType, t reflect.Type, funcName string, structName string) (err error) {
	errMsg := ""
	defer func() {
		if errMsg != "" {
			err = errors.New(fmt.Sprintf("%s.%s 方法 %s", structName, funcName, errMsg))
		}
	}()
	if t.Kind() != reflect.Func {
		errMsg = "必须是func类型"
		return
	}
	funcT := constants.FuncTypeMap[ft]

	ctxt := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !t.In(0).Implements(ctxt) {
		errMsg = fmt.Sprintf("第一个参数必须为context")
	}
	if !t.In(1).ConvertibleTo(reflect.TypeOf(&gorm.DB{})) {
		errMsg = fmt.Sprintf("第二个参数必须为*gorm.db")
	}
	errt := reflect.TypeOf((*error)(nil)).Elem()
	if !t.Out(t.NumOut() - 1).Implements(errt) {
		errMsg = fmt.Sprintf("返回值最后一个参数必须为err")
	}
	needInCount := len(funcT.In) + 2
	needOutCount := len(funcT.Out) + 1
	if t.NumIn() < needInCount {
		errMsg = fmt.Sprintf("参数数量不对哦.need: %d, actual: %d", needInCount, t.NumIn())
	}
	if t.NumOut() != needOutCount {
		errMsg = fmt.Sprintf("返回值数量不对哦.need: %d, actual: %d", needOutCount, t.NumOut())
	}
	return nil
}

func fixSelOp(sels []*filter.Where, t reflect.Type) {
	if len(sels) == 1 {
		sels[0].Essential = true
	}
	for i, sel := range sels {
		inType := t.In(t.NumIn() - 1 - i)
		if strings.ToLower(sel.Op) == "=" &&
			(inType.Kind() == reflect.Array || inType.Kind() == reflect.Slice) {
			sel.Op = "in"
		}
	}
}

func CheckFuncType(ft constants.FuncType, t reflect.Type, sels []*filter.Where, funcName string, structName string) (err error) {
	errMsg := ""
	defer func() {
		if errMsg != "" {
			err = errors.New(fmt.Sprintf("%s.%s 方法 %s", structName, funcName, errMsg))
		}
	}()
	funcT := constants.FuncTypeMap[ft]
	for i, v := range funcT.In {
		switch v {
		case constants.ModelTypeIns:
			if t.In(i+2).Kind() != reflect.Ptr {
				errMsg = fmt.Sprintf("第%d个参数必须为model指针", i+3)
			}
		case constants.ArgsTypeIns:
			if len(sels) == 0 {
				errMsg = fmt.Sprintf("函数为BY函数但是没有Selector")
			}
			skip := 0
			for i, sel := range sels {
				if sel.HasSetObj() {
					skip++
					continue
				}
				inType := t.In(t.NumIn() - 1 + skip - i)
				switch strings.ToLower(sel.Op) {
				case "in":
					if inType.Kind() != reflect.Array && inType.Kind() != reflect.Slice {
						errMsg = fmt.Sprintf("select %s in 对应参数不是数组", sel.Col)
					}
				case "=", "<", ">", "<>", "!=", ">=", "<=":
					switch inType.Kind() {
					case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8,
						reflect.Bool, reflect.String:
					default:
						errMsg = fmt.Sprintf("select %s %s 对应参数不是数字或字符串或bool值", sel.Col, sel.Op)
					}
				}
			}
		}
	}
	for i, v := range funcT.Out {
		switch v {
		case constants.ModelArrTypeIns:
			arrType := t.Out(i).Kind()
			if arrType != reflect.Slice && arrType != reflect.Array {
				errMsg = fmt.Sprintf("第%d个返回值必须为model数组", i+1)
			}
		case constants.ModelTypeIns:
			arrType := t.Out(i).Kind()
			if arrType != reflect.Ptr || t.Out(i).Elem().Kind() != reflect.Struct {
				errMsg = fmt.Sprintf("第%d个返回值必须为model结构体指针", i+1)
			}
		}
	}
	return
}

func FillFunc(dbFunc *finisher.DBFunc) (err error) {
	funcType := parser.GetFuncType(dbFunc)
	dbFunc.FuncType = funcType
	filters := parser.ParseFilter(dbFunc)
	dbFunc.Filters = filters
	finisher := MakeFinisher(funcType, dbFunc)
	dbFunc.Finisher = finisher
	err = dbFunc.Check()
	if err != nil {
		return err
	}
	f := reflect.MakeFunc(dbFunc.T, func(args []reflect.Value) (results []reflect.Value) {
		results, err := dbFunc.Exec(args)
		if err != nil {
			results = append(results, reflect.ValueOf(err))
		} else {
			results = append(results, constants.NilError)
		}
		return results
	})
	dbFunc.V.Set(f)
	return nil
}

func MakeFinisher(funcType constants.FuncType, dbFunc *finisher.DBFunc) (res finisher.DBFinisher) {
	switch funcType {
	case constants.Get, constants.GetBy:
		t := dbFunc.T.Out(0).Elem()
		return &finisher.Find{
			First: true,
			ResBO: dbFunc.IsBOT(t),
		}
	case constants.MGetBy, constants.MGet, constants.GetByCondition:
		t := dbFunc.T.Out(0)
		return &finisher.Find{
			First: false,
			ResBO: dbFunc.IsBOT(t),
		}
	case constants.Update:
		up := utils.GetTagStr(dbFunc.Tag, "update", "u")
		if up == "" {
			return &finisher.Update{
				UpdateModel: true,
			}
		}
		return parser.WalkUpdaters(up, dbFunc)
	case constants.Create:
		return &finisher.Create{}
	case constants.Count:
		return &finisher.Count{}
	case constants.Delete:
		up := utils.GetTagStr(dbFunc.Tag, "delete", "d", "del")
		if up == "" {
			return &finisher.Delete{}
		}
		return parser.WalkUpdaters(up, dbFunc)
	}
	return res
}

func RegisterWriteDB(f func(context.Context) *gorm.DB) {
	finisher.GetWriteDB = f
}

func RegisterReadDB(f func(context.Context) *gorm.DB) {
	finisher.GetReadDB = f
}

func Transaction(ctx context.Context, db *gorm.DB, f func(tx *gorm.DB) error) (err error) {
	if db == nil && finisher.GetWriteDB == nil {
		return errors.New("gorm-curd未定义GetWriteDB方法, 请使用curd.RegisterWriteDB定义")
	}
	if db == nil {
		db = finisher.GetWriteDB(ctx)
	}
	return db.Transaction(f)
}

func GetTime(millsSeconds int64) time.Time {
	if millsSeconds == 0 {
		return time.Time{}
	}
	return time.Unix(0, millsSeconds*int64(time.Millisecond))
}

func GetJSONString(object interface{}) string {
	if object == nil {
		return ""
	}
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return "{}" //default empty json object string
	}
	return string(jsonBytes)
}
