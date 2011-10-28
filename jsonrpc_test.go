package jsonrpc

import (
	"testing"
	"reflect"
	"fmt"
)

type MyInt int

func (this MyInt) PublicMethod1(a, b *int) int {
	return *a + *b
}

func (this MyInt) PublicMethod2() {}
func (this MyInt) PublicMethod3() {}
func (this MyInt) PublicMethod4() {}

func (this MyInt) privateMethod1() {}
func (this MyInt) privateMethod2() {}

func TestCall(t *testing.T) {
	a, b, o := 4, 5, MyInt(0)
	call := &Call {
		MethodName: "PublicMethod1",
		Parameters: []interface{} {
			a, b,
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

func TestWrongCall(t *testing.T) {
	a, o := 4, MyInt(0)
	call := &Call {
		MethodName: "PublicMethod1",
		Parameters: []interface{} {
			a,
		},
	}

	rpc := New(o)
	_, e := rpc.ExecuteCall(call)
	if e != ErrNumArguments {
		t.Fatalf("Call does not work: %s", e.String())
	}
}

func TestEnumerate(t *testing.T) {
	a, b, o := 4, 5, MyInt(0)
	call := &Call {
		MethodName: "_enumerate",
		Parameters: []interface{} {
			a, b,
		},
	}

	rpc := New(o)
	r, e := rpc.ExecuteCall(call)
	if e != nil {
		t.Fatalf("_enmerate does not work: %s", e.String())
	}
	if len(r) != 4 {
		t.Fatalf("_enumerate did not return the right amount of functions, namely %d", len(r))
	}
	for i, m := range r {
		r2, ok := m.(Method)
		if !ok {
			t.Fatalf("_enumerate did not return a method enumeration at index %d, but a %s", i, reflect.ValueOf(r[0]).Type().String())
		}
		expected := fmt.Sprintf("PublicMethod%d", i+1)
		if r2.Name != expected {
			t.Fatalf("_enumerate did not return the \"%s\" but \"%s\"", expected, r2.Name)
		}
	}
}

func TestMarshalling(t *testing.T) {
	o := MyInt(0)
	call := `{
		"MethodName": "PublicMethod1",
		"Parameters": [
			1, 2
		]
	}`

	rpc := New(o)
	r, e := rpc.Execute(call)
	if e != nil {
		t.Fatalf("Call does not work: %s", e.String())
	}
	if r != "[3]" {
		t.Fatalf("Call did return the right result: \"%s\"", r)
	}
}
