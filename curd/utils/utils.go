package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"sync"

	"gorm.io/gorm"
)

type Pagination struct {
	PageNo     int64
	PageSize   int64
	TotalCount int64
	HasMore    bool
}

func NilPtr(t reflect.Type) reflect.Value {
	return reflect.Indirect(reflect.New(t))
}

func IsParamsValid(params ...interface{}) (res bool) {
	res = true
	for _, param := range params {
		switch v := reflect.ValueOf(param); v.Kind() {
		case reflect.String:
			res = res && v.String() != ""
			if !res {
				return
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			res = res && v.Int() != 0
			if !res {
				return
			}
		case reflect.Ptr:
			res = res && !v.IsNil()
			if !res {
				return
			}
		case reflect.Invalid:
			res = false
			if !res {
				return
			}
		case reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
			res = res && v.Len() != 0
			if !res {
				return
			}
		default:
			res = res && !reflect.ValueOf(param).IsZero()
			if !res {
				return
			}
		}
	}
	return res
}

func Paginate(pageInfo *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pageInfo.PageNo
		if page == 0 {
			page = 1
		}

		pageSize := pageInfo.PageSize
		switch {
		case pageSize > 1000:
			pageSize = 1000
		case pageSize <= 0:
			pageSize = 1
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

func GetJSONStr(obj interface{}) string {
	str, _ := json.Marshal(obj)
	return string(str)
}

var (
	smap sync.Map
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func ToDBName(name string) string {
	if name == "" {
		return ""
	} else if v, ok := smap.Load(name); ok {
		return fmt.Sprint(v)
	}

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool // upper case == true
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase
	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}

	return buf.String()
}

func GetTagStr(tag reflect.StructTag, args ...string) string {
	if args == nil || len(args) == 0 {
		return ""
	}
	var res string
	for _, arg := range args {
		res = tag.Get(arg)
		if res != "" {
			return res
		}
	}
	return ""
}

func Atoi(s string) int {
	res, _ := strconv.Atoi(s)
	return res
}
