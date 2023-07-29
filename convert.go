package govaluate

import (
	"encoding/json"
	"strconv"
)

func convert2Str(value interface{}) string {

	if value == nil {
		return ""
	}

	switch value.(type) {
	case string:
		return value.(string)
	case float64:
		s := value.(float64)
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		s := value.(float32)
		return strconv.FormatFloat(float64(s), 'f', -1, 64)
	case uint:
		s := value.(uint)
		return strconv.Itoa(int(s))
	case uint8:
		s := value.(uint8)
		return strconv.Itoa(int(s))
	case uint16:
		s := value.(uint16)
		return strconv.Itoa(int(s))
	case uint32:
		s := value.(uint32)
		return strconv.Itoa(int(s))
	case uint64:
		s := value.(uint64)
		return strconv.Itoa(int(s))
	case int:
		s := value.(int)
		return strconv.Itoa(int(s))
	case int8:
		s := value.(int8)
		return strconv.Itoa(int(s))
	case int16:
		s := value.(int16)
		return strconv.Itoa(int(s))
	case int32:
		s := value.(int32)
		return strconv.Itoa(int(s))
	case int64:
		s := value.(int64)
		return strconv.Itoa(int(s))
	case []byte:
		s := value.([]byte)
		return string(s)
	default:
		s, _ := json.Marshal(value)
		return string(s)
	}
}

func convert2Float64(value interface{}) (float64, error) {

	switch value.(type) {
	case float64:
		return value.(float64), nil
	case float32:
		s := value.(float32)
		return float64(s), nil
	case uint:
		s := value.(uint)
		return float64(s), nil
	case uint8:
		s := value.(uint8)
		return float64(s), nil
	case uint16:
		s := value.(uint16)
		return float64(s), nil
	case uint32:
		s := value.(uint32)
		return float64(s), nil
	case uint64:
		s := value.(uint64)
		return float64(s), nil
	case int:
		s := value.(int)
		return float64(s), nil
	case int8:
		s := value.(int8)
		return float64(s), nil
	case int16:
		s := value.(int16)
		return float64(s), nil
	case int32:
		s := value.(int32)
		return float64(s), nil
	case int64:
		s := value.(int64)
		return float64(s), nil
	default:
		s := convert2Str(value)
		return strconv.ParseFloat(s, 64)
	}
}
