package parser

import (
	"reflect"
	"regexp"
	"strings"

	finisher2 "github.com/JingruiLea/gobatis/curd/oprator/finisher"

	"github.com/JingruiLea/gobatis/curd/constants"

	"github.com/JingruiLea/gobatis/curd/utils"

	"github.com/JingruiLea/gobatis/curd/oprator/filter"
)

func ParseFilter(dbFunc *finisher2.DBFunc) (res []filter.DBFilter) {
	res = make([]filter.DBFilter, 0)
	tagS := string(dbFunc.Tag)
	tagS = strings.TrimSpace(tagS)
	tagS = strings.ReplaceAll(tagS, "\n", " ")
	tagS = strings.ReplaceAll(tagS, "\t", " ")
	for strings.Contains(tagS, "  ") {
		tagS = strings.ReplaceAll(tagS, "  ", " ")
	}
	tagS = strings.ReplaceAll(tagS, ": ", ":")
	tagS = strings.ReplaceAll(tagS, " :", ":")
	tagS = strings.ReplaceAll(tagS, "\" ", "\"")
	tagS = strings.ReplaceAll(tagS, " \"", "\"")
	dbFunc.Tag = reflect.StructTag(tagS)
	tableName := utils.ToDBName(strings.TrimRight(dbFunc.StructName, "Dal"))
	modelTagStr := utils.GetTagStr(dbFunc.Tag, "table")
	if modelTagStr == "" {
		model := &filter.Model{
			Table: tableName,
		}
		model.Model = dbFunc.ModelT
		res = append(res, model)
	}
	tag := dbFunc.Tag
	tagStr := string(tag)
	re, _ := regexp.Compile("[a-z^:]*:[^:]*\"")
	tagArr := re.FindAllString(tagStr, -1)
	selTagStr := utils.GetTagStr(dbFunc.Tag, "s", "sel", "select")
	funcNameSelStr := parseFuncName(dbFunc.Name)
	var sels = make([]*filter.Where, 0)
	if selTagStr == "" && funcNameSelStr != "" {
		sels = WalkSel(funcNameSelStr, dbFunc)
		for _, sel := range sels {
			res = append(res, sel)
		}
	}
	for _, s := range tagArr {
		headers := strings.Split(s, ":\"")
		h := headers[0]
		data := strings.TrimRight(headers[1], "\"")
		switch h {
		case "count":
			if data == "?" {
				res = append(res, &filter.Count{
					DefaultFilter: filter.DefaultFilter{Active: true},
				})
			} else if data == "res" {
				res = append(res, &filter.Count{AsRes: true})
			}
		case "dis", "distinct":
			if data == "?" {
				res = append(res, &filter.Distinct{DefaultFilter: filter.DefaultFilter{Active: true}})
			} else {
				res = append(res, &filter.Distinct{V: data})
			}
		case "g", "group":
			if data == "?" {
				res = append(res, &filter.Group{DefaultFilter: filter.DefaultFilter{Active: true}})
			} else {
				res = append(res, &filter.Group{V: data})
			}
		case "l", "limit":
			if data == "?" {
				res = append(res, &filter.Limit{DefaultFilter: filter.DefaultFilter{Active: true}})
			} else {
				res = append(res, &filter.Limit{V: utils.Atoi(data)})
			}
		case "lock":
			if data == "?" {
				res = append(res, &filter.Lock{DefaultFilter: filter.DefaultFilter{Active: true}})
			} else {
				res = append(res, &filter.Lock{V: true})
			}
		case "offset":
			if data == "?" {
				res = append(res, &filter.Offset{DefaultFilter: filter.DefaultFilter{Active: true}})
			} else {
				res = append(res, &filter.Offset{V: utils.Atoi(data)})
			}
		case "order":
			if data == "?" {
				res = append(res, &filter.Order{DefaultFilter: filter.DefaultFilter{Active: true}})
			} else {
				res = append(res, &filter.Order{V: data})
			}
		case "p", "page":
			res = append(res, &filter.Page{
				DefaultFilter: filter.DefaultFilter{Active: true},
			})
		case "s", "sel", "select":
			subsels := WalkSel(data, dbFunc)
			for _, sel := range subsels {
				res = append(res, sel)
			}
			sels = append(sels, subsels...)
		}
	}
	if len(sels) == 0 && dbFunc.FuncType != constants.Create && utils.GetTagStr(dbFunc.Tag, "s", "sel", "select") != "ALL" {
		res = append(res, &filter.Where{
			DefaultFilter: filter.DefaultFilter{Active: true},
		})
	}
	return res
}
