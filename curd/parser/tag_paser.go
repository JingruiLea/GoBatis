package parser

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/JingruiLea/gobatis/curd/oprator/finisher"
	"google.golang.org/appengine/log"

	"github.com/JingruiLea/gobatis/curd/oprator/filter"

	"github.com/JingruiLea/gobatis/curd/constants"
	"github.com/JingruiLea/gobatis/curd/utils"
)

func GetFuncType(dbFunc *finisher.DBFunc) constants.FuncType {
	funcName := utils.GetTagStr(dbFunc.Tag, "method", "m")
	if funcName == "" {
		funcName = strings.ToLower(dbFunc.Name)
	}
	if strings.HasPrefix(funcName, string(constants.Get)) {
		if strings.HasSuffix(funcName, constants.ConditionSuffix) {
			return constants.GetByCondition
		}
		return constants.Get
	}
	if strings.HasPrefix(funcName, string(constants.MGet)) {
		return constants.MGet
	}
	if strings.HasPrefix(funcName, string(constants.Create)) {
		return constants.Create
	}
	if strings.HasPrefix(funcName, string(constants.Delete)) {
		return constants.Delete
	}
	if strings.HasPrefix(funcName, string(constants.Update)) {
		return constants.Update
	}
	countStr := utils.GetTagStr(dbFunc.Tag, "count", "c")
	if countStr != "" {
		return constants.Count
	}
	return ""
}

func GetSelectStr(f *finisher.DBFunc) string {
	selStr := parseTag(f.Tag)
	if selStr == "" {
		selStr = parseFuncName(f.Name)
	}
	return selStr
}

func parseTag(t reflect.StructTag) string {
	selector, ok := t.Lookup("s")
	if !ok {
		selector = t.Get("select")
	}
	return selector
}

func parseFuncName(f string) string {
	fArr := strings.Split(f, "By")
	if len(fArr) < 2 {
		return ""
	}
	if strings.HasPrefix(f, "M") {
		return utils.ToDBName(fArr[len(fArr)-1]) + " in"
	}
	return utils.ToDBName(fArr[len(fArr)-1])
}

func WalkUpdaters(updateStr string, dbFunc *finisher.DBFunc) (res *finisher.Update) {
	updateArr := strings.Split(updateStr, ",")
	res = &finisher.Update{
		Updates: make(map[string]interface{}),
	}
	for _, ups := range updateArr {
		uarr := strings.Split(strings.TrimSpace(ups), " ")
		switch len(uarr) {
		case 0:
			return res
		case 1, 2:
			res.Set(uarr[0], nil)
			return
		case 3:
			if uarr[2] != "?" {
				res.Set(uarr[0], getStrValue(uarr[2]))
			} else {
				res.Set(uarr[0], nil)
			}
		default:
			upskv := strings.SplitN(strings.TrimSpace(ups), "=", 2)
			if len(upskv) < 2 {
				return
			}
			exprString := strings.TrimSpace(upskv[1])
			expr := finisher.Expr{}
			for _, char := range exprString {
				if char == '?' {
					expr.Count++
				}
			}
			expr.V = exprString
			res.Set(strings.TrimSpace(upskv[0]), expr)
		}
	}
	return
}

func WalkSel(str string, dbFunc *finisher.DBFunc) []*filter.Where {
	var m = make(utils.Queue, 0)
	sels := make([]*filter.Where, 0)
	if str == "" {
		return sels
	}
	strArr := strings.Split(str, " ")
	m.Append("and")
	for i := 0; i < len(strArr); i++ {
		switch strArr[i] {
		case "and", "or":
			sels = append(sels, subWalk(&m, dbFunc))
		}
		m.Append(strArr[i])
	}
	sels = append(sels, subWalk(&m, dbFunc))
	return sels
}

func subWalk(s *utils.Queue, dbFunc *finisher.DBFunc) *filter.Where {
	res := &filter.Where{Op: "=", Or: false, Essential: true, Obj: "?"}
	if s.Length() > 4 {
		log.Errorf("%s.%s select tag parse error! ", dbFunc.StructName, dbFunc.Name)
	}
	if word, _ := s.Pop(); word == "or" {
		res.Or = true
	}
	word, ok := s.Pop()
	if ok {
		if strings.HasSuffix(word, "?") {
			res.Essential = false
			word = strings.TrimRight(word, "?")
		}
		res.Col = utils.ToDBName(word)
	}
	word, ok = s.Pop()
	if ok {
		word = strings.ToLower(word)
		if IsOp(word) {
			res.Op = word
		}
	}
	word, ok = s.Pop()
	if ok {
		res.Obj = getStrValue(word)
	}
	if !res.HasSetObj() {
		res.Active = true
	}
	return res
}

func getStrValue(word string) interface{} {
	if strings.HasPrefix(word, "'") && strings.HasPrefix(word, "'") {
		return word[1 : len(word)-1]
	}
	num, err := strconv.Atoi(word)
	if err != nil {
		return word
	}
	return num
}

func IsOp(s string) bool {
	for _, s2 := range constants.OpArr {
		if s2 == s {
			return true
		}
	}
	return false
}

func GetValue(word string) (res string, intValue int64) {
	if strings.HasPrefix(word, "'") && strings.HasPrefix(word, "'") {
		return word[1 : len(word)-1], 0
	}
	num, err := strconv.ParseInt(word, 10, 64)
	if err != nil {
		return word, 0
	}
	return "", num
}
