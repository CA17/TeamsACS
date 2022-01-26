package vread

import (
	"errors"
	"testing"
)

func TestStringReader_ReadRequiedString(t *testing.T) {
	var getStr = func()(string, error) {
		return "", errors.New("err")
	}
	var getStr2 = func()(string, error) {
		return "2", nil
	}
	var v, v2 string
	NewStringReader().ReadRequiedString(&v, *NewSP(getStr()))
	t.Log(v)
	NewStringReader().ReadRequiedString(&v2, *NewSP(getStr2()))
	t.Log(v2)
}

func TestStringReader_ReadRequiedString2(t *testing.T) {
	vmap := map[string]interface{}{
		"v1":1,
		"v2":"2",
	}
	var v1, v2 string
	sr := NewStringReader().
		ReadRequiedString2(&v1, vmap, "v1").
		ReadRequiedString2(&v2, vmap,"v2")
	if sr.LastError!= nil {
		t.Fatal(sr.LastError)
	}
	t.Log(v1, v2)
}

func TestStringReader_ReadRequiedString3(t *testing.T) {
	var getStr = func()(string, error) {
		return "1", nil;
	}
	var getStr2 = func()(string, error) {
		return "2", nil
	}
	var v1, v2 string
	sr := NewStringReader().
		ReadRequiedString(&v1, *NewSP(getStr())).
		ReadRequiedString(&v2, *NewSP(getStr2()))
	if sr.LastError!= nil {
		t.Fatal(sr.LastError)
	}
	t.Log(v1, v2)
}
