// JsonRPC makes it possible to invode
// public methods of objects with a JSON formatted
// protocol.
// A limitation is, that parameters have to be pointer types
// (or interfaces)
package jsonrpc

import (
	"reflect"
	"json"
	"os"
)

// Calls will be umarshalled into this struct.
type Call struct {
	MethodName string
	Parameters []interface{}
}

type JsonRPC struct {
	object  reflect.Value
	methods map[string]reflect.Method
}

// Creates a new JsonRPC object
// which delegates calls to given object’s
// public methods.
func New(obj interface{}) *JsonRPC {
	rpc := &JsonRPC{
		object: reflect.ValueOf(obj),
	}
	rpc.cacheMethods()
	return rpc
}

func (this *JsonRPC) Execute(marshalled_call string) ([]interface{}, os.Error) {
	call := &Call{}
	e := json.Unmarshal([]byte(marshalled_call), call)
	if e != nil {
		return nil, e
	}
	return this.ExecuteCall(call)
}

var (
	ErrNoSuchMethod = os.NewError("Method does not exist")
)

func (this *JsonRPC) ExecuteCall(call *Call) ([]interface{}, os.Error) {
	// Catch the special function "_enumerate" which lists all
	// available methods
	if call.MethodName == "_enumerate" {
		return this.getMethodsAsInterface(), nil
	}

	method, ok := this.methods[call.MethodName]
	if !ok {
		return nil, ErrNoSuchMethod
	}
	return executeCall(this.object, method, call.Parameters), nil
}

type Method struct {
	Name      string
	NumParams int
}

func (this *JsonRPC) getMethodsAsInterface() (m []interface{}) {
	methods := this.getMethods()
	m = make([]interface{}, len(methods))
	for i := range methods {
		m[i] = methods[i]
	}
	return m
}

func (this *JsonRPC) getMethods() (m []Method) {
	m = make([]Method, this.object.NumMethod())
	for i := 0; i < this.object.NumMethod(); i++ {
		method := this.object.Type().Method(i)
		// first parameter is the receiver object
		// which is provided by us, not the remote caller
		m[i].Name = method.Name
		m[i].NumParams = this.object.Method(i).Type().NumIn() - 1
		i++
	}
	return
}
