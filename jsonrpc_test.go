package jsonrpc

import (
	"testing"
	"reflect"
)

type MyInt int

func (this MyInt) PublicMethod1(a, b *int) int {
	return *a + *b
}

func TestCall(t *testing.T) {
	a, b, o := 4, 5, MyInt(0)
	call := &Call {
		MethodName: "PublicMethod1",
		Parameters: []interface{} {
			&a, &b,
		},
	}

	rpc := New(o)
	r, e := rpc.ExecuteCall(call)
	if e != nil {
		t.Fatalf("Call does not work: %s", e.String())
	}
	if len(r) != 1 {
		t.Fatalf("Call did return more results than expected, namely %d", len(r))
	}
	r3, ok := r[0].(int)
	if !ok {
		t.Fatalf("Call did not return an int: %s", reflect.ValueOf(r[0]).Type().String())
	}
	if r3 != 9 {
		t.Fatalf("Call did not yield 9 but %d", r3)
	}
}

func TestEnumerate(t *testing.T) {
	a, b, o := 4, 5, MyInt(0)
	call := &Call {
		MethodName: "_enumerate",
		Parameters: []interface{} {
			&a, &b,
		},
	}

	rpc := New(o)
	r, e := rpc.ExecuteCall(call)
	if e != nil {
		t.Fatalf("_enmerate does not work: %s", e.String())
	}
	if len(r) != 1 {
		t.Fatalf("_enumerate did not return the right amount of functions, namely %d", len(r))
	}
	r2, ok := r[0].(Method)
	if !ok {
		t.Fatalf("_enumerate did not return a method enumeration", reflect.ValueOf(r[0]).Type().String())
	}
	if r2.Name != "PublicMethod1" {
		t.Fatalf("_enumerate did not return the right function name: %s", r2.Name)
	}
}
