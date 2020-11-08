package filter

import (
	"reflect"

	"github.com/JingruiLea/gobatis/curd/utils"
	"gorm.io/gorm"
)

type Page struct {
	Pagination interface{}
	DefaultFilter
}

func (p *Page) SetValue(v reflect.Value) {
	p.Pagination = v.Interface()
}

func (p *Page) Filter(db *gorm.DB) (*gorm.DB, error) {
	pagePtr := reflect.ValueOf(p.Pagination)
	if pagePtr.IsNil() {
		return db, nil
	}
	pageInfoV := pagePtr.Elem()
	pageNo := pageInfoV.FieldByName("PageNo")
	pageSize := pageInfoV.FieldByName("PageSize")
	myPageInfo := &utils.Pagination{
		PageNo:   pageNo.Int(),
		PageSize: pageSize.Int(),
	}
	err := db.Count(&myPageInfo.TotalCount).Error
	pageInfoV.FieldByName("TotalCount").Set(reflect.ValueOf(myPageInfo.TotalCount))
	if myPageInfo.PageSize*myPageInfo.PageNo < myPageInfo.TotalCount {
		pageInfoV.FieldByName("HasMore").Set(reflect.ValueOf(true))
	}
	var res *gorm.DB
	res = db.Scopes(utils.Paginate(myPageInfo))
	return res, err
}

func (p *Page) Check(in []reflect.Type, i int) (string, int) {
	t := in[i]
	if t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Struct {
		return "page对应的参数必须为结构体指针", i + 1
	}
	t = t.Elem()
	pn, ok := t.FieldByName("PageNo")
	if !ok || !(pn.Type.Kind() == reflect.Int || pn.Type.Kind() == reflect.Int64 || pn.Type.Kind() == reflect.Int32) {
		return "Page没有合法的PageNo", i + 1
	}
	ps, ok := t.FieldByName("PageSize")
	if !ok || !(ps.Type.Kind() == reflect.Int || ps.Type.Kind() == reflect.Int64 || ps.Type.Kind() == reflect.Int32) {
		return "Page没有合法的PageSize", i + 1
	}
	tc, ok := t.FieldByName("TotalCount")
	if !ok || !(tc.Type.Kind() == reflect.Int || tc.Type.Kind() == reflect.Int64 || tc.Type.Kind() == reflect.Int32) {
		return "Page没有合法的TotalCount", i + 1
	}
	hm, ok := t.FieldByName("HasMore")
	if !ok || !(hm.Type.Kind() == reflect.Bool) {
		return "Page没有合法的HasMore", i + 1
	}
	return "", i + 1
}
