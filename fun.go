package winterMVC

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const patternProjectUrl = `^v\d+$`

// MustCompileVersion path 为 version 则返回 true
func MustCompileVersion(path string) bool {
	reg := regexp.MustCompile(patternProjectUrl)
	return reg.MatchString(path)
}

// 将string参数转为typeStr指定类型的值
// typeStr string	类型字串	支持int、float、bool、Time
// valueStr string	值
// interface{}	为 nil则失败
func stringToType(typeStr string, valueStr string) interface{} {
	var data interface{}
	var e error
	switch typeStr {
	case "int":
		if "" == valueStr {
			return 0
		}
		data, e = strconv.Atoi(valueStr)
	case "uint":
		if "" == valueStr {
			return uint(0)
		}
		val, e1 := strconv.Atoi(valueStr)
		if nil == e1 {
			data = uint(val)
		} else {
			e = e1
			data = val
		}
	case "int8":
		if "" == valueStr {
			return int8(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 8)
		if nil == e {
			data = int8(data.(int64))
		}
	case "uint8":
		if "" == valueStr {
			return uint8(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 8)
		if nil == e {
			data = uint8(data.(uint64))
		}
	case "int16":
		if "" == valueStr {
			return int16(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 16)
		if nil == e {
			data = int16(data.(int64))
		}
	case "uint16":
		if "" == valueStr {
			return uint16(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 16)
		if nil == e {
			data = uint16(data.(uint64))
		}
	case "int32":
		if "" == valueStr {
			return int32(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 32)
		if nil == e {
			data = int32(data.(int64))
		}
	case "uint32":
		if "" == valueStr {
			return uint32(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 32)
		if nil == e {
			data = uint32(data.(uint64))
		}
	case "int64":
		if "" == valueStr {
			return int64(0)
		}
		data, e = strconv.ParseInt(valueStr, 10, 64)
	case "uint64":
		if "" == valueStr {
			return uint64(0)
		}
		data, e = strconv.ParseUint(valueStr, 10, 64)
	case "bool":
		if "" == valueStr {
			return false
		}
		data, e = strconv.ParseBool(valueStr)
	case "float32":
		if "" == valueStr {
			return 0
		}
		data, e = strconv.ParseFloat(valueStr, 32)
		if nil == e {
			data = float32(data.(float64))
		}
	case "float64":
		if "" == valueStr {
			return 0
		}
		data, e = strconv.ParseFloat(valueStr, 64)
	case "string":
		if "" == valueStr {
			return ""
		}
		data = valueStr
	case "Time":
		if 10 == len(valueStr) {
			data, e = time.Parse("2006-01-02", valueStr)
		} else if 13 == len(valueStr) {
			data, e = time.Parse("2006-01-02 15", valueStr)
		} else if 16 == len(valueStr) {
			data, e = time.Parse("2006-01-02 15:04", valueStr)
		} else if 19 == len(valueStr) {
			data, e = time.Parse("2006-01-02 15:04:05", valueStr)
		} else {
			//data, e = time.Parse("2006-01-02'T'15:04:05.999 Z", valueStr)
		}
		if nil != e {
			e = nil
			data = time.Now()
		}
	}
	if nil != e {
		logError("stringToType", e)
		return nil
	}
	return data
}

// 字串数组转字串，以,拼接
// strArr []string	字串数组
// string	以【】间隔的值
func stringArrayToString(strArr []string) string {
	str := ""
	for inx, _ := range strArr {
		if 0 == inx {
			str = strArr[inx]
			continue
		}
		str = fmt.Sprintf("%s【】%s", str, strArr[inx])
	}
	return str
}
