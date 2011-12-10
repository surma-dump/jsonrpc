package jsonrpc

import (
	"reflect"
	"unicode"
	"unicode/utf8"
	"encoding/json"
)

// Fuck yeah, UTF8 compatible!
func isPublicMethod(method reflect.Method) bool {
	return unicode.IsUpper(utf8.NewString(method.Name).At(0))

}

func (this *JsonRPC) cacheMethods() {
	this.methods = make(map[string]reflect.Method)
	num := this.object.Type().NumMethod()
	for i := 0; i < num; i++ {
		method := this.object.Type().Method(i)
		if isPublicMethod(method) {
			this.methods[method.Name] = method
		}
	}
}

func createUnmarshallingInterface(pointerType reflect.Type) interface{} {
	if pointerType.Kind() != reflect.Ptr {
		panic("Parameters have to be pointers")
	}

	valueType := pointerType.Elem()
	valuePointer := reflect.New(valueType)
	return valuePointer.Interface()
}

func createParameterArray(method reflect.Method) []interface{} {
	numIn := method.Type.NumIn()
	r := make([]interface{}, numIn-1)
	// First parameter is receiver, no unmarshalling
	// necessary there, hence i := 1
	for i := 1; i < method.Type.NumIn(); i++ {
		intype := method.Type.In(i)
		r[i-1] = createUnmarshallingInterface(intype)
	}
	return r
}

func typeParams(method reflect.Method, raw_params []interface{}) []reflect.Value {
	params := createParameterArray(method)

	data, e := json.Marshal(raw_params)
	if e != nil {
		panic(e)
	}

	e = json.Unmarshal(data, &params)
	if e != nil {
		panic(e)
	}

	result := make([]reflect.Value, len(params))
	for i := range params {
		result[i] = reflect.ValueOf(params[i])
	}
	return result
}

func executeCall(rcv reflect.Value, method reflect.Method, raw_params []interface{}) []interface{} {
	params := []reflect.Value{rcv}
	params = append(params, typeParams(method, raw_params)...)
	reflected_results := method.Func.Call(params)
	result := value2Interface(reflected_results)
	return result
}

func interface2Value(in []interface{}) []reflect.Value {
	out := make([]reflect.Value, 0)
	for _, i := range in {
		out = append(out, reflect.ValueOf(i))
	}
	return out
}

func value2Interface(in []reflect.Value) []interface{} {
	out := make([]interface{}, 0)
	for _, i := range in {
		out = append(out, i.Interface())
	}
	return out
}
