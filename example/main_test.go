package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/JingruiLea/gobatis/curd/utils"

	"gorm.io/gorm"

	"github.com/JingruiLea/gobatis/example/model"

	"github.com/JingruiLea/gobatis/curd"
)

var ctx context.Context

func init() {
	curd.RegisterReadDB(GetDB)
	curd.RegisterWriteDB(GetDB)
	err := curd.AutoFill(&UserDalIns)
	if err != nil {
		panic(err)
		return
	}
	ctx = context.Background()
}

//###
func Test1(t *testing.T) {
	err := UserDalIns.Create(ctx, nil, &model.User{
		ID:        0,
		UserID:    12345,
		Username:  "李井瑞",
		Age:       2,
		Phone:     "1234",
		IsDeleted: false,
	})
	if err != nil {
		panic(err)
	}
}

//###
func Test2(t *testing.T) {
	_, err := UserDalIns.GetUser(ctx, nil, 10086)
	if err != nil {
		panic(err)
	}
}

func Test3(t *testing.T) {
	_, err := UserDalIns.GetUserWithLock(ctx, nil, 10086, true)
	if err != nil {
		panic(err)
	}
	_, err = UserDalIns.GetUserWithLock(ctx, nil, 10086, false)
	if err != nil {
		panic(err)
	}
}

//###
func Test4(t *testing.T) {
	_, err := UserDalIns.MGetUser(ctx, nil, []int64{1, 2, 3}, []string{"hell0", "1243"}, 0)
	if err != nil {
		panic(err)
	}
	_, err = UserDalIns.MGetUser(ctx, nil, []int64{}, []string{"hell0", "1243"}, 0)
	if err != nil {
		panic(err)
	}
	_, err = UserDalIns.MGetUser(ctx, nil, []int64{1, 2, 3}, []string{}, 50)
	if err != nil {
		panic(err)
	}
}

func Test5(t *testing.T) {
	err := UserDalIns.UpdateUserAge(ctx, nil, 123, 10, "hello")
	if err != nil {
		panic(err)
	}
	err = UserDalIns.UpdateUserAge(ctx, nil, 123, 0, "cc")
	if err != nil {
		panic(err)
	}
	err = UserDalIns.UpdateUserAge(ctx, nil, 123, 0, "")
	if err != nil {
		panic(err)
	}
}

func Test6(t *testing.T) {
	err := UserDalIns.DeleteUser(ctx, nil, 123)
	if err != nil {
		panic(err)
	}
	err = UserDalIns.DeleteUserV2(ctx, nil, 123)
	if err != nil {
		panic(err)
	}
}

func Test7(t *testing.T) {
	err := UserDalIns.UpdateUserAgeV2(ctx, nil, 123)
	if err != nil {
		panic(err)
	}
	err = UserDalIns.UpdateUserAgeV3(ctx, nil, 123, 2)
	if err != nil {
		panic(err)
	}
}

func Test8(t *testing.T) {
	err := curd.Transaction(ctx, nil, func(tx *gorm.DB) error {
		_, err := UserDalIns.GetUserWithLock(ctx, tx, 10086, true)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Test9(t *testing.T) {
	page := &utils.Pagination{
		PageNo:     1, // 改为1
		PageSize:   4,
		TotalCount: 0,
		HasMore:    false,
	}
	users, err := UserDalIns.MGetUserPage(ctx, nil, page, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println(utils.GetJSONStr(page))
	fmt.Println(utils.GetJSONStr(users))
}

func Test10(t *testing.T) {
	UserDalIns.CreateBO(ctx, nil, &model.BO{})
}

func Test11(t *testing.T) {
	//curd.FillForDebug(UserDalIns, "GetBO")
	res, _ := UserDalIns.GetBO(ctx, nil, 12345)
	fmt.Println(res)
	res2, _ := UserDalIns.MGetBO(ctx, nil, 12345)
	fmt.Println(res2)

}

func Test12(t *testing.T) {
	UserDalIns.UpdateBO(ctx, nil, 123, &model.BO{})
}

func Test13(t *testing.T) {
	ctx := context.Background()
	db := GetDB(ctx)
	res := &model.User{}
	ret := db.Model(&model.User{}).Where(`user_id > ?`, 0).Take(res)
	fmt.Println(ret.RowsAffected)
	if err := ret.Error; err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}

func Test14(t *testing.T) {
	ctx := context.Background()
	db := GetDB(ctx)
	err := db.Create(&model.User{
		ID:        0,
		UserID:    0,
		Username:  "",
		Age:       0,
		Phone:     "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		IsDeleted: false,
	}).Error
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Test15(t *testing.T) {
	ctx := context.WithValue(context.Background(), "123", "")
	db := GetDB(ctx)
	db.Transaction(func(tx *gorm.DB) error {
		err := db.Updates(&model.User{
			ID:        10086,
			UserID:    0,
			Username:  "",
			Age:       0,
			Phone:     "",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			IsDeleted: false,
		}).Error
		if err != nil {
			fmt.Println(err.Error())
		}
		return err
	})

}
