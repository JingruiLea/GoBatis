package constants

import (
	"context"
	"errors"
	"reflect"

	"github.com/JingruiLea/gobatis/curd/utils"

	"gorm.io/gorm"
)

var OpArr = []string{
	"in", ">", "<", "<>", "!=", "=", ">=", "<=",
}

var NilError = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
var ctxType context.Context
var gormDBType *gorm.DB

type FuncType string

type DalFunc func(ctx context.Context, db *gorm.DB, funcT reflect.Type, args []reflect.Value) (results []reflect.Value, err error)

var DalMap = map[FuncType]DalFunc{}

const (
	Get   FuncType = "get"
	GetBy FuncType = "getby"

	MGet   FuncType = "mget"
	MGetBy FuncType = "mgetby"

	Update   FuncType = "update"
	UpdateBy FuncType = "updateby"

	Delete   FuncType = "del"
	DeleteBy FuncType = "deleteby"

	Create FuncType = "create"

	GetByCondition FuncType = "getbycondition"

	ConditionSuffix string = "bycondition"

	Count FuncType = "count"
)

type Model struct {
}

var modelType *Model

var ModelTypeIns = reflect.TypeOf(modelType)

var modelArrType []*Model
var ModelArrTypeIns = reflect.TypeOf(modelArrType)

var ErrTypeIns = reflect.TypeOf(errors.New(""))

type args struct {
}

var ArgsTypeIns = reflect.TypeOf(&args{})

type FuncT struct {
	In  []reflect.Type
	Out []reflect.Type
}

var FuncTypeMap = map[FuncType]FuncT{
	//func(context.Context, *gorm.DB, *model.User)(*model.User,error)
	Get: {
		In: []reflect.Type{
			ModelTypeIns,
		},
		Out: []reflect.Type{
			ModelTypeIns,
		},
	},
	//func(context.Context, *gorm.DB, *model.User)([]*model.User,error)
	MGet: {
		In: []reflect.Type{
			ModelTypeIns,
		},
		Out: []reflect.Type{
			ModelArrTypeIns,
		},
	},
	//func(context.Context, *gorm.DB, *model.User)(error)
	Update: {
		In: []reflect.Type{
			ModelTypeIns,
		},
	},
	//func(context.Context, *gorm.DB, *model.User)(error)
	Delete: {
		In: []reflect.Type{
			ModelTypeIns,
		},
	},
	//func(context.Context, *gorm.DB, args)(*model.User, error)
	GetBy: {
		In: []reflect.Type{
			ArgsTypeIns,
		},
		Out: []reflect.Type{
			ModelTypeIns,
		},
	},
	//func(context.Context, *gorm.DB, args)([]*model.User, error)
	MGetBy: {
		In: []reflect.Type{
			ArgsTypeIns,
		},
		Out: []reflect.Type{
			ModelArrTypeIns,
		},
	},
	UpdateBy: {
		In: []reflect.Type{
			ArgsTypeIns,
		},
		Out: []reflect.Type{
			ModelArrTypeIns,
		},
	},
	DeleteBy: {
		In: []reflect.Type{
			ArgsTypeIns,
		},
		Out: []reflect.Type{
			ModelArrTypeIns,
		},
	},
	Create: {
		In: []reflect.Type{
			ModelTypeIns,
		},
		Out: []reflect.Type{},
	},
	GetByCondition: {
		In: []reflect.Type{
			reflect.TypeOf(&utils.Pagination{}),
			ArgsTypeIns,
		},
		Out: []reflect.Type{ModelArrTypeIns},
	},
}
