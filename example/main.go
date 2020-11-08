package main

import (
	"context"

	"github.com/JingruiLea/gobatis/curd/utils"

	"github.com/JingruiLea/gobatis/example/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var UserDalIns = UserDal{}

type UserDal struct {
	model    *model.User
	bo       *model.BO
	Create   func(ctx context.Context, db *gorm.DB, user *model.User) error
	CreateBO func(ctx context.Context, db *gorm.DB, user *model.BO) error

	GetBO    func(ctx context.Context, db *gorm.DB, uid int64) (*model.BO, error)   `select:"user_id = ?"`
	MGetBO   func(ctx context.Context, db *gorm.DB, uid int64) ([]*model.BO, error) `select:"user_id = ?"`
	UpdateBO func(ctx context.Context, db *gorm.DB, uid int64, d *model.BO) error   `select:"user_id = ?"`

	Update func(ctx context.Context, db *gorm.DB, user *model.User, d *model.User) error

	GetUserWithLock func(ctx context.Context, db *gorm.DB, uid int64, lock bool) (*model.User, error) `select:"user_id = ?" lock:"?"`
	GetUser         func(ctx context.Context, db *gorm.DB, uid int64) (*model.User, error)            `select:"user_id = ?"`

	MGetUser func(ctx context.Context, db *gorm.DB, uid []int64, phone []string, age int) ([]*model.User, error) `s:"user_id? in ? or phone? in ? and age? = ?"`

	UpdateUserAge func(ctx context.Context, db *gorm.DB, uid int64, age int64, name string) error `s:"user_id" update:"age? = ?, username? = ?"`
	UpdateUser    func(ctx context.Context, db *gorm.DB, user *model.User, age int) error         `u:"age = ?"`

	MGetUserV3   func(ctx context.Context, db *gorm.DB, age int64) ([]*model.User, error)                         `s:"age > ?"`
	MGetUserPage func(ctx context.Context, db *gorm.DB, page *utils.Pagination, age int64) ([]*model.User, error) `page:"?" s:"age > ?"`

	UpdateUserPhoneByUserID func(ctx context.Context, db *gorm.DB, uid int64, phone string, age int) (int64, error) `s:"user_id" u:"phone? = ?, age? = ?"`
	UpdateUserAgeV2         func(ctx context.Context, db *gorm.DB, uid int64) error                                 `s:"user_id" u:"age = age + 1"`
	UpdateUserAgeV3         func(ctx context.Context, db *gorm.DB, uid int64, param1 int) error                     `s:"user_id" u:"age = age + ?"`

	UpdateNoneTag   func(ctx context.Context, db *gorm.DB, uid int64, user *model.User) error `s:"user_id"`
	UpdateNoneTagV2 func(ctx context.Context, db *gorm.DB, selectModel *model.User, updateModel *model.User) error
	DeleteUser      func(ctx context.Context, db *gorm.DB, uid int64) error                              `s:"user_id"`
	DeleteUserV2    func(ctx context.Context, db *gorm.DB, uid int64) error                              `s:"user_id" delete:"is_deleted = 1"`
	MGetUserOlder   func(ctx context.Context, db *gorm.DB, count *int64, age int) ([]*model.User, error) `count:"?" s:"age > ?" `

	MGetUserAndCountOlder      func(ctx context.Context, db *gorm.DB, count *int64, age int) ([]*model.User, error) `count:"?" s:"age > ?"`
	MGetValidUserAndCountOlder func(ctx context.Context, db *gorm.DB, age int, count *int64) ([]*model.User, error) `s:"age > ?" count:"?" s:"is_deleted? = 0"`
}

func (UserDal) RecordNotFound(ctx context.Context, structName string, funcName string) error {

	return nil
}

func GetDB(ctx context.Context) *gorm.DB {
	dsn := "root:admin@tcp(127.0.0.1:3306)/baozhao?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})
	return db.Debug().WithContext(ctx)
}

func main() {

}
