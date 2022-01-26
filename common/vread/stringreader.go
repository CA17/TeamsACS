package vread

import (
	"errors"
)

// 参数快速读取,避免多次错误处理
type StringReader struct {
	LastError error
}

func NewStringReader() *StringReader {
	return &StringReader{}
}

type StringParam struct {
	value string
	err   error
}

func NewSP(value string, err error) *StringParam {
	return &StringParam{value: value, err: err}
}

func (sr *StringReader) ReadRequiedString(ref *string, p StringParam) *StringReader {
	if sr.LastError != nil {
		return sr
	}
	if p.err != nil {
		sr.LastError = p.err
	} else {
		*ref = p.value
	}
	return sr
}

func (sr *StringReader) ReadRequiedString2(ref *string, vmap map[string]interface{}, name string) *StringReader {
	if sr.LastError != nil {
		return sr
	}
	v, ok := vmap[name]
	if !ok || v == "" {
		sr.LastError = errors.New(name + " is not set")
		return sr
	}
	vstr, ok := v.(string)
	if !ok {
		sr.LastError = errors.New(name + " is not string")
		return sr
	}
	*ref = vstr
	return sr
}
