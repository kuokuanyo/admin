package modules

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
)

func InArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// 耞把计del肚del+s(把计)+del┪[s(把计)]
func Delimiter(del, s string) string {
	if del == "[" {
		return "[" + s + "]"
	}
	return del + s + del
}

func FilterField(filed, delimiter string) string {
	if delimiter == "[" {
		return "[" + filed + "]"
	}
	return delimiter + filed + delimiter
}

func InArrayWithoutEmpty(arr []string, str string) bool {
	if len(arr) == 0 {
		return true
	}
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func RemoveBlankFromArray(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func Uuid() string {
	uid, _ := uuid.NewV4()
	rst := uid.String()
	return rst
}

func SetDefault(source, def string) string {
	if source == "" {
		return def
	}
	return source
}

func GetPage(page string) (pageInt int) {
	if page == "" {
		pageInt = 1
	} else {
		pageInt, _ = strconv.Atoi(page)
	}
	return
}

func AorB(condition bool, a string, b string) string {
	if condition {
		return a
	}
	return b
}

func AorEmpty(condition bool, a string) string {
	if condition {
		return a
	}
	return ""
}
