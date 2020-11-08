package finisher

import (
	"reflect"

	"gorm.io/gorm"
)

type DBFinisher interface {
	Finish(*gorm.DB, *DBFunc) ([]reflect.Value, error)
	Fill([]reflect.Value)
	CheckRes(o []reflect.Type) string
}
