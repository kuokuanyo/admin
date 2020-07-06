// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
)

// SetColVarType set the column type.
// 設置欄位類型
func SetColVarType(colVar *[]interface{}, i int, typeName string) {
	//判別欄位數值類別
	switch {
	//Contains 在types.go中
	//contains中的第二個參數(欄位的類型)�m在 types.go中
	// 將字串(typeName)轉換成 type DatabaseType
	case Contains(DT(typeName), BoolTypeList):
		var s sql.NullBool
		(*colVar)[i] = &s
	case Contains(DT(typeName), IntTypeList):
		var s sql.NullInt64
		(*colVar)[i] = &s
	case Contains(DT(typeName), FloatTypeList):
		var s sql.NullFloat64
		(*colVar)[i] = &s
	case Contains(DT(typeName), UintTypeList):
		var s []uint8
		(*colVar)[i] = &s
	case Contains(DT(typeName), StringTypeList):
		var s sql.NullString
		(*colVar)[i] = &s
	default:
		var s interface{}
		(*colVar)[i] = &s
	}
}

// SetResultValue set the result value.
// 設置數值
func SetResultValue(result *map[string]interface{}, index string, colVar interface{}, typeName string) {
	switch {
	case Contains(DT(typeName), BoolTypeList):
		temp := *(colVar.(*sql.NullBool))
		if temp.Valid {
			(*result)[index] = temp.Bool
		} else {
			(*result)[index] = nil
		}
	case Contains(DT(typeName), IntTypeList):
		temp := *(colVar.(*sql.NullInt64))
		if temp.Valid {
			(*result)[index] = temp.Int64
		} else {
			(*result)[index] = nil
		}
	case Contains(DT(typeName), FloatTypeList):
		temp := *(colVar.(*sql.NullFloat64))
		if temp.Valid {
			(*result)[index] = temp.Float64
		} else {
			(*result)[index] = nil
		}
	case Contains(DT(typeName), UintTypeList):
		(*result)[index] = *(colVar.(*[]uint8))
	case Contains(DT(typeName), StringTypeList):
		temp := *(colVar.(*sql.NullString))
		if temp.Valid {
			(*result)[index] = temp.String
		} else {
			(*result)[index] = nil
		}
	default:
		if colVar2, ok := colVar.(*interface{}); ok {
			if colVar, ok = (*colVar2).(int64); ok {
				(*result)[index] = colVar
			} else if colVar, ok = (*colVar2).(string); ok {
				(*result)[index] = colVar
			} else if colVar, ok = (*colVar2).(float64); ok {
				(*result)[index] = colVar
			} else if colVar, ok = (*colVar2).([]uint8); ok {
				(*result)[index] = colVar
			} else {
				(*result)[index] = colVar
			}
		}
	}
}
