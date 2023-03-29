package qgen

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertToEscapeString(obj interface{}, def string) (res string) {
	switch v := obj.(type) {
	case int, int64, int32, float64, float32, bool:
		res = fmt.Sprintf("%v", v)
	case time.Time:
		res = (obj.(time.Time)).Format("2006-01-02 15:04:05")
	case string:
		switch v {
		case "__jsonNOW()__":
			res = "DATE_FORMAT(NOW(), '%Y-%m-%dT%TZ')"
		case "__NOW()__":
			res = "NOW()"
		default:
			res = strconv.Quote(fmt.Sprintf("%v", v))
		}
	case []string:
		res = "( "
		data := obj.([]string)
		for idx, v2 := range data {
			res += strconv.Quote(fmt.Sprintf("%s", v2))
			if idx < len(data)-1 {
				res += ", "
			} else {
				res += " )"
			}
		}
	case []int:
		res = "( "
		data := obj.([]int)
		for idx, v2 := range data {
			res += fmt.Sprintf("%d", v2)
			if idx < len(data)-1 {
				res += ", "
			} else {
				res += " )"
			}
		}
	case []int64:
		res = "( "
		data := obj.([]int64)
		for idx, v2 := range data {
			res += fmt.Sprintf("%d", v2)
			if idx < len(data)-1 {
				res += ", "
			} else {
				res += " )"
			}
		}
	case []float64:
		res = "( "
		data := obj.([]float64)
		for idx, v2 := range data {
			res += fmt.Sprintf("%f", v2)
			if idx < len(data)-1 {
				fmt.Printf(", ")
			} else {
				fmt.Printf(" )")
			}
		}
	default:
		res = def
	}

	return
}
