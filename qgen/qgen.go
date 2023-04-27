package qgen

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/DeniesKresna/gohelper/utinterface"
	"github.com/DeniesKresna/gohelper/utlog"
	"github.com/DeniesKresna/gohelper/utslice"
)

type Args struct {
	Offset     int64
	Limit      int
	Sorting    []string
	Conditions map[string]interface{}
	Fields     []string
	Groups     []string
	Distinct   bool
	Updates    map[string]interface{}
}

type Obj struct {
	IsLogged        bool
	ListTableColumn map[string]map[string]string
}

func InitObject(isLogged bool, tables ...interface{}) (obj *Obj, err error) {
	obj = &Obj{
		IsLogged: isLogged,
	}

	var listTableColumn = make(map[string]map[string]string)

	for _, tbSt := range tables {
		var tbVal = reflect.ValueOf(tbSt)
		if !utinterface.IsStruct(tbVal) {
			err = errors.New("Table should be struct")
			return
		}

		var tbName string
		tbNameRes := tbVal.MethodByName("GetTableName").Call([]reflect.Value{})
		tbName = tbNameRes[0].Interface().(string)

		if tbName == "" {
			err = errors.New("Error when get table name")
			return
		}

		listTableColumn[tbName] = make(map[string]string)
		listTableColumn[tbName]["*"] = "*"

		var recursiveStructField func(structField interface{}) error
		recursiveStructField = func(structField interface{}) (errs error) {
			var tbV = reflect.ValueOf(structField)
			var reflectType = tbV.Type()

			for i := 0; i < reflectType.NumField(); i++ {
				if reflectType.Field(i).Type.Kind() == reflect.Struct {
					errs = recursiveStructField(tbV.Field(i).Interface())
					if errs != nil {
						return
					}
				}

				fieldTags := reflectType.Field(i).Tag

				dbTagStr := fieldTags.Get("db")
				if dbTagStr == "" {
					continue
				}
				dbTags := strings.Split(dbTagStr, ",")
				dbTag := dbTags[0]

				var jsonDBTag string
				jsonDBTagStr := fieldTags.Get("jsondb")
				if jsonDBTagStr == "" {
					jsonDBTag = dbTag
				} else {
					jsonDBTags := strings.Split(jsonDBTagStr, ",")
					jsonDBTag = jsonDBTags[0]
				}

				if dbTag != "" && jsonDBTag != "" {
					if _, ok := listTableColumn[tbName][dbTag]; ok {
						errs = errors.New(fmt.Sprintf("listTableColumn[%s][%s] has been exist", tbName, dbTag))
						return
					}
					listTableColumn[tbName][dbTag] = jsonDBTag
				}
			}
			return
		}

		err = recursiveStructField(tbVal.Interface())
		if err != nil {
			return
		}
	}
	obj.ListTableColumn = listTableColumn
	return
}

func (q *Obj) standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func (q *Obj) Build(query string, args Args) (res string) {
	var format = make(map[string]interface{})
	query = q.standardizeSpaces(query)

	err := json.Unmarshal([]byte(query), &format)
	if err != nil {
		utlog.Errorf("err: %+v\n", err)
		return
	}

	// Resolve Conditions
	var condsCol = make(map[string]string)
	{
		for k, v := range args.Conditions {
			var wv string
			condsCols := strings.Split(k, ":")
			if len(condsCols) == 1 {
				wv = fmt.Sprintf("= %s", ConvertToEscapeStringSQL(v, ""))
			} else if len(condsCols) == 2 {
				wv = fmt.Sprintf("%s %s", condsCols[1], ConvertToEscapeStringSQL(v, ""))
			}

			condsCol[condsCols[0]] = wv
		}
	}

	buildType := "query"
	if len(args.Updates) > 0 {
		buildType = "update"
	}

	res = q.RecursiveBuild(format, buildType, args, condsCol, nil)

	if q.IsLogged {
		utlog.Infof("QUERY LOG: %s\n", res)
	}

	return
}

func (q *Obj) RecursiveBuild(form interface{}, kind string, args Args, condsCol map[string]string, selectAlias map[string]string) (res string) {
	switch kind {
	case "query":
		format, ok := form.(map[string]interface{})
		if !ok {
			utlog.Warnf("Cannot parse %s\n", kind)
			return
		}

		var (
			subString   = ""
			fromAlias   = make(map[string]string)
			joinAlias   = make(map[string]string)
			joinConn    = make(map[string]string)
			selectAlias = make(map[string]string)
		)

		// Handle FROM
		{
			fromVal, fromExist := format["from"]
			if !fromExist {
				utlog.Warn("Select should be has from\n")
				return
			}

			fromValMap, ok := fromVal.(map[string]interface{})
			if !ok {
				utlog.Warnf("Cannot parse %+v\n", fromVal)
				return
			}

			fromValue, ok := fromValMap["value"]
			if !ok {
				utlog.Warn("from should has col")
				return
			}

			var as string
			fromAs, ok := fromValMap["as"]
			if !ok {
				utlog.Warn("from should has as")
				return
			}
			as, ok = fromAs.(string)
			if !ok {
				utlog.Warn("as should be string")
				return
			}

			fromValueMap, ok := fromValue.(map[string]interface{})
			if ok {
				rs := q.RecursiveBuild(fromValueMap, "query", args, condsCol, nil)
				fromAlias[as] = rs
			} else {
				fromValueStr, ok := fromValue.(string)
				if ok {
					fromAlias[as] = fromValueStr
				} else {
					utlog.Warn("from should be has map or string value")
					return
				}
			}
		}

		// Handler JOIN
		{
			joinVal, joinExist := format["join"]
			if joinExist {
				joinSlice, ok := joinVal.([]interface{})
				if !ok {
					utlog.Warnf("Join should be slice")
					return
				}

				for _, join := range joinSlice {
					joinMap, ok := join.(map[string]interface{})
					if !ok {
						utlog.Warn("join slice data should be map")
						return
					}

					var as string
					joinAs, ok := joinMap["as"]
					if !ok {
						utlog.Warn("join slice data should has as")
						return
					}
					as, ok = joinAs.(string)
					if !ok {
						utlog.Warn("as should be string")
						return
					}

					joinVal, ok := joinMap["value"]
					if !ok {
						utlog.Warn("join slice data should has value")
						return
					}

					var joinValueStr string
					joinValMap, ok := joinVal.(map[string]interface{})
					if ok {
						rs := q.RecursiveBuild(joinValMap, "query", args, condsCol, nil)
						joinValueStr = fmt.Sprintf("(%s)", rs)
					} else {
						joinValueStr, ok = joinVal.(string)
						if !ok {
							utlog.Warn("join should be has map or string value")
							return
						}
					}

					joinType, ok := joinMap["type"]
					if !ok {
						utlog.Warn("join slice data should has type")
						return
					}

					joinTypeStr, ok := joinType.(string)
					if !ok {
						utlog.Warn("join type should be string")
						return
					}

					joinCon, ok := joinMap["conn"]
					if !ok {
						utlog.Warn("join slice data should has conn")
						return
					}

					joinConStr, ok := joinCon.(string)
					if !ok {
						utlog.Warn("join conn should be string")
						return
					}

					if _, ok := joinAlias[as]; !ok {
						joinAlias[as] = joinValueStr
					}

					joinConn[as] = fmt.Sprintf("%s join %s %s on %s", joinTypeStr, joinValueStr, as, joinConStr)
				}
			}
		}

		// Handler SELECT
		{
			selectVal, selectExist := format["select"]
			if !selectExist {
				utlog.Warn("Query should has select\n")
				return
			}

			subString = "select"

			if args.Distinct == true {
				subString += " distinct"
			}

			selectSlice, ok := selectVal.([]interface{})
			if !ok {
				utlog.Warn("Select should be slice")
				return
			}

			for _, sel := range selectSlice {
				selMap, ok := sel.(map[string]interface{})
				if !ok {
					utlog.Warn("Couldnt read select data")
					continue
				}

				selCol, ok := selMap["col"]
				if !ok {
					utlog.Warn("Select should has col")
					continue
				}

				selColStr, ok := selCol.(string)
				if !ok {
					utlog.Warn("Select column should be string")
					continue
				}

				if selColStr != "-" {
					selColStr = strings.TrimSpace(selColStr)
					selColStrs := strings.Split(selColStr, ".")
					selAlias, selField := selColStrs[0], selColStrs[1]
					tableName, ok := fromAlias[selAlias]
					if !ok {
						tableName, ok = joinAlias[selAlias]
						if !ok {
							continue
						}
					}

					if selField == "*" {
						for key, v := range q.ListTableColumn[tableName] {
							xField := fmt.Sprintf("%s.%s", selAlias, key)
							selectAlias[xField] = v

							if !utslice.IsExist(args.Fields, xField) {
								continue
							}
							var sc string
							if key != v {
								jsonPaths := strings.Split(v, ">")
								switch len(jsonPaths) {
								case 2, 3:
									sc = fmt.Sprintf("%s.%s->>\"%s\"", selAlias, jsonPaths[0], jsonPaths[1])
								default:
									sc = xField
								}
							} else {
								sc = xField
							}
							if strings.HasSuffix(sc, "*") {
								subString += fmt.Sprintf(" %s,", sc)
								continue
							}
							subString += fmt.Sprintf(" %s as %s,", sc, key)
						}
						continue
					}

					if !utslice.IsExist(args.Fields, selColStr) {
						continue
					}
				}

				selVal, ok := selMap["value"]
				if !ok {
					utlog.Warn("Select should has value except * format")
					continue
				}

				selValStr, ok := selVal.(string)
				if !ok {
					utlog.Warn("Select value should be string")
					continue
				}

				var selAsStr string
				selAs, ok := selMap["as"]
				if ok {
					selAsStr, ok = selAs.(string)
					if !ok {
						utlog.Warn("Select as should be string")
						continue
					}

					selAsStr = fmt.Sprintf(" as %s", selAsStr)
				}

				subString += fmt.Sprintf(" %s%s,", selValStr, selAsStr)
			}
		}

		// Add FROM
		for k, v := range fromAlias {
			subString = RemoveLastCommas(subString)
			subString += fmt.Sprintf(" from %s %s", v, k)
		}

		// Add JOIN
		for _, v := range joinConn {
			subString += fmt.Sprintf(" %s", v)
		}

		// Handler WHERE
		{
			whereVal, whereExist := format["where"]
			if whereExist {
				subString += " where "
				var sc string
				whereValMap, ok := whereVal.(map[string]interface{})
				if !ok {
					utlog.Warn("Cannot parse Where\n")
					return
				}
				whereAnd, ok := whereValMap["and"]
				if ok {
					sc = q.RecursiveBuild(whereAnd, "and", args, condsCol, selectAlias)
				} else {
					whereOr, ok := whereValMap["or"]
					if ok {
						sc = q.RecursiveBuild(whereOr, "or", args, condsCol, selectAlias)
					} else {
						utlog.Warn("Where should has and or or key\n")
						return
					}
				}

				subString += fmt.Sprintf("( %s )", sc)
			}
		}

		// Handler LIMIT, OFFSET, SORT
		{
			if len(args.Sorting) > 0 {
				sortRes := "order by"
				for idx, v := range args.Sorting {
					isDesc := false
					if strings.HasPrefix(v, "-") {
						isDesc = true
					}
					newV := strings.TrimPrefix(v, "-")

					ordColStr := strings.TrimSpace(newV)
					ordColStrs := strings.Split(ordColStr, ".")
					ordTbAlias, ordField := ordColStrs[0], ordColStrs[1]
					tableName, ok := fromAlias[ordTbAlias]
					if !ok {
						tableName, ok = joinAlias[ordTbAlias]
						if !ok {
							continue
						}
					}
					valueCol := q.ListTableColumn[tableName][ordField]
					valueCols := strings.Split(valueCol, ">")
					switch len(valueCols) {
					case 2, 3:
						newV = fmt.Sprintf("%s.%s->>\"%s\"", ordTbAlias, valueCols[0], valueCols[1])
					default:
					}
					if ok {
						sortRes += fmt.Sprintf(" %s ", newV)
						if isDesc {
							sortRes += "desc"
						} else {
							sortRes += "asc"
						}
						if idx < len(args.Sorting)-1 {
							sortRes += ","
						}
					}
				}
				subString += fmt.Sprintf(" %s", sortRes)
			}

			if args.Limit > 0 {
				subString += fmt.Sprintf(" limit %d", args.Limit)
			}

			if args.Offset > 0 {
				subString += fmt.Sprintf(" offset %d", args.Offset)
			}
		}

		res = subString

		return

	case "update":
		format, ok := form.(map[string]interface{})
		if !ok {
			utlog.Warnf("Cannot parse %s\n", kind)
			return
		}

		var (
			fromAlias   = make(map[string]string)
			selectAlias = make(map[string]string)
		)

		// Handle FROM
		{
			fromVal, fromExist := format["from"]
			if !fromExist {
				utlog.Warn("Update should be has from\n")
				return
			}

			fromValMap, ok := fromVal.(map[string]interface{})
			if !ok {
				utlog.Warnf("Cannot parse %+v\n", fromVal)
				return
			}

			fromValue, ok := fromValMap["value"]
			if !ok {
				utlog.Warn("from should has col")
				return
			}

			var as string
			fromAs, ok := fromValMap["as"]
			if !ok {
				utlog.Warn("from should has as")
				return
			}
			as, ok = fromAs.(string)
			if !ok {
				utlog.Warn("as should be string")
				return
			}

			fromValueMap, ok := fromValue.(map[string]interface{})
			if ok {
				rs := q.RecursiveBuild(fromValueMap, "query", args, condsCol, nil)
				fromAlias[as] = rs
			} else {
				fromValueStr, ok := fromValue.(string)
				if ok {
					fromAlias[as] = fromValueStr
				} else {
					utlog.Warn("from should be has map or string value")
					return
				}
			}
		}

		subString := "update "

		for k, v := range fromAlias {
			subString += fmt.Sprintf(" %s %s", v, k)
		}

		subString += " set"

		// Handler SELECT
		{
			selectVal, selectExist := format["set"]
			if !selectExist {
				utlog.Warn("Update should has set")
				return
			}

			selectSlice, ok := selectVal.([]interface{})
			if !ok {
				utlog.Warn("Set should be slice")
				return
			}

			for _, sel := range selectSlice {
				selMap, ok := sel.(map[string]interface{})
				if !ok {
					utlog.Warn("Couldnt read set data")
					continue
				}

				selCol, ok := selMap["col"]
				if !ok {
					utlog.Warn("Set should has col")
					continue
				}

				selColStr, ok := selCol.(string)
				if !ok {
					utlog.Warn("Set column should be string")
					continue
				}

				if selColStr != "-" {
					selColStr = strings.TrimSpace(selColStr)
					selColStrs := strings.Split(selColStr, ".")
					selAlias, selField := selColStrs[0], selColStrs[1]
					tableName, ok := fromAlias[selAlias]
					if !ok {
						continue
					}

					if selField == "*" {
						for key, v := range q.ListTableColumn[tableName] {
							xField := fmt.Sprintf("%s.%s", selAlias, key)
							selectAlias[xField] = v

							updateFieldVal, oks := args.Updates[xField]
							if !oks {
								continue
							}
							var sc string
							if key != v {
								jsonPaths := strings.Split(v, ">")
								switch len(jsonPaths) {
								case 2, 3:
									sc = fmt.Sprintf(" %s = JSON_SET(%s, '%s', %s),", jsonPaths[0], jsonPaths[0], jsonPaths[1], ConvertToEscapeStringSQL(updateFieldVal, ""))
								default:
									sc = fmt.Sprintf(" %s = %s,", key, ConvertToEscapeString(updateFieldVal, ""))
								}
							} else {
								sc = fmt.Sprintf(" %s = %s,", key, ConvertToEscapeString(updateFieldVal, ""))
							}
							subString += sc
						}
						continue
					}
					continue
				}

				selVal, ok := selMap["value"]
				if !ok {
					utlog.Warn("Set should has value except * format")
					continue
				}

				selValStr, ok := selVal.(string)
				if !ok {
					utlog.Warn("Set value should be string")
					continue
				}

				selCond, ok := selMap["update_value"]
				if ok {
					selCondStr, ok := selCond.(string)
					if !ok {
						utlog.Warn("Set Condition should be string")
						continue
					}

					subString += fmt.Sprintf(" %s = %s,", selValStr, ConvertToEscapeString(selCondStr, ""))
					continue
				}
				subString += fmt.Sprintf(" %s, ", selValStr)
				continue
			}
		}

		subString = RemoveLastCommas(subString)

		// Handler WHERE
		{
			whereVal, whereExist := format["where"]
			if whereExist {
				subString += " where "
				var sc string
				whereValMap, ok := whereVal.(map[string]interface{})
				if !ok {
					utlog.Warn("Cannot parse Where\n")
					return
				}
				whereAnd, ok := whereValMap["and"]
				if ok {
					sc = q.RecursiveBuild(whereAnd, "and", args, condsCol, selectAlias)
				} else {
					whereOr, ok := whereValMap["or"]
					if ok {
						sc = q.RecursiveBuild(whereOr, "or", args, condsCol, selectAlias)
					} else {
						utlog.Warn("Where should has and or or key\n")
						return
					}
				}

				subString += fmt.Sprintf("( %s )", sc)
			}
		}

		res = subString
	case "and", "or":
		var (
			sc string
		)

		format, ok := form.(map[string]interface{})
		if ok {
			_, ok = format[kind]
			if ok {
				sc = q.RecursiveBuild(form, kind, args, condsCol, nil)
				sc += fmt.Sprintf("( %s ) %s ", sc, kind)
			}
		} else {
			format2Slice, ok := form.([]interface{})
			if !ok {
				utlog.Warn("And or Or should has map or array of data")
				return
			}

			for _, whereRow := range format2Slice {
				whereRowMap, ok := whereRow.(map[string]interface{})
				if !ok {
					utlog.Warn("Cannot parse whereRowMap")
					return
				}

				whereCol, ok := whereRowMap["col"]
				if !ok {
					utlog.Warn("whereRowMap should has col key")
					return
				}

				whereColStr, ok := whereCol.(string)
				if !ok {
					utlog.Warn("Where Col should be string")
					return
				}

				var (
					conValue    string
					isCondExist bool
				)
				if whereColStr != "-" {
					conValue, isCondExist = condsCol[whereColStr]
					if !isCondExist {
						continue
					}
				}

				whereVal, ok := whereRowMap["value"]
				if !ok {
					utlog.Warn("whereRowMap should has value key")
					return
				}

				whereValStr, ok := whereVal.(string)
				if ok {
					valueVal, ok := selectAlias[whereValStr]
					if ok {
						whereValStr = strings.TrimSpace(whereValStr)
						whereValStrs := strings.Split(whereValStr, ".")
						selAlias := whereValStrs[0]

						valueCols := strings.Split(valueVal, ">")
						switch len(valueCols) {
						case 3:
							valueVal = fmt.Sprintf("JSON_VALUE(%s.%s, '%s' RETURNING %s)", selAlias, valueCols[0], valueCols[1], valueCols[2])
						case 2:
							valueVal = fmt.Sprintf("%s.%s->>\"%s\"", selAlias, valueCols[0], valueCols[1])
						default:
							valueVal = whereValStr
						}
					} else {
						valueVal = whereValStr
					}
					sc += fmt.Sprintf(" %s %s %s", valueVal, conValue, kind)
				} else {
					whereValMap, ok := whereVal.(map[string]interface{})
					if ok {
						scx := q.RecursiveBuild(whereValMap, "query", args, condsCol, nil)
						sc += fmt.Sprintf(" ( %s ) %s %s", scx, conValue, kind)
					}
				}
			}
		}

		sc = strings.TrimSpace(sc)

		if strings.HasSuffix(sc, kind) {
			sc = sc[:len(sc)-len(kind)]
		}

		res = sc
		return
	default:

	}
	return
}
