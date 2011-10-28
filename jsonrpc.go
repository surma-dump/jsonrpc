// JsonRPC makes it possible to invoke
// public methods of objects with a JSON formatted
// protocol.
// A limitation is, that parameters have to be pointer types
// (or interfaces).
// A special function "_enumerate" with no parameters
// returns a list of callable functions.
package jsonrpc

import (
	"reflect"
	"json"
	"os"
	"utf8"
	"unicode"
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

// Executes the call described by the passed json string.
// This is basically unmarshalles the string into a Call struct
// and calls ExecuteCall() afterwards
// The return parameters of the called functions are save into the array.
func (this *JsonRPC) Execute(marshalled_call string) (string, os.Error) {
	call := &Call{}
	e := json.Unmarshal([]byte(marshalled_call), call)
	if e != nil {
		return "", e
	}
	returns, e := this.ExecuteCall(call)
	if e != nil {
		return "", e
	}
	marshalled_returns, e := json.Marshal(returns)
	if e != nil {
		return "", e
	}
	return string(marshalled_returns), nil
}

var (
	ErrNoSuchMethod = os.NewError("Method does not exist")
	ErrNumArguments = os.NewError("Wrong number of arguments")
)

// Executes the call described by the given call struct.
// The return parameters of the called functions are saved into the array.
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

	if method.Type.NumIn() - 1 != len(call.Parameters) {
		return nil, ErrNumArguments
	}
	return executeCall(this.object, method, call.Parameters), nil
}

type Method struct {
	Name      string
	NumParams int
}

func (this *JsonRPC) getMethodsAsInterface() (m []interface{}) {
	methods := this.GetMethods()
	m = make([]interface{}, len(methods))
	for i := range methods {
		m[i] = methods[i]
	}
	return m
}

// Returns all available methods
func (this *JsonRPC) GetMethods() (m []Method) {
	m = make([]Method, this.object.NumMethod())
	j := 0
	for i := 0; i < this.object.NumMethod(); i++ {
		method := this.object.Type().Method(i)
		// Is this a public method?
		if unicode.IsUpper(utf8.NewString(method.Name).At(0)) {
			m[j].Name = method.Name
			// first parameter is the receiver object
			// which is provided by us, not the remote caller
			m[j].NumParams = this.object.Method(i).Type().NumIn() - 1
			j++
		}
	}
	return m[0:j]
}
