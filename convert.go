package utils

import (
	"fmt"
	"strconv"
)

type Str string

func (s Str) Exist() bool {
	return string(s) != string(0x1E)
}

func (s Str) String() string {
	if s.Exist() {
		return string(s)
	}
	return ""
}

func (s Str) Unit8() (uint8, error) {
	v, err := strconv.ParseUint(s.String(), 10, 8)
	return uint8(v), err
}

func (s Str) Int() (int, error) {
	v, err := strconv.ParseInt(s.String(), 10, 0)
	return int(v), err
}

func (s Str) Int64() (int64, error) {
	v, err := strconv.ParseInt(s.String(), 10, 64)
	return int64(v), err
}
func (s Str) MustUint8() uint8 {
	v, _ := s.Unit8()
	return v
}

func (s Str) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s Str) MustInt64() int64 {
	v, _ := s.Int64()
	return v
}

type argInt []int

func (a argInt) Get(i int, arg ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(arg) > 0 {
		r = arg[0]
	}
	return
}

func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}
