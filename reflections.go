package jsrpc

import (
	"pusher"
	"reflect"
	"unicode"
	"utf8"
	"json"
)

func parseInterface(v interface{}) (c callback) {
	c.object = reflect.ValueOf(v)
	c.methods = make(map[string]reflect.Method)
	num := c.object.Type().NumMethod()
	for i := 0; i < num; i++ {
		method := c.object.Type().Method(i)
		if isPublicMethod(method) {
			c.methods[method.Name] = method
		}
	}
	return
}

// Fuck yeah, UTF8 compatible!
func isPublicMethod(method reflect.Method) bool {
	return unicode.IsUpper(utf8.NewString(method.Name).At(0))

}

func jsonify(v interface{}) string {
	b, e := json.Marshal(v)
	if e != nil {
		panic("Marshalling failed: " + e.String())
	}
	return string(b)
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

func decodeParams(reflectParams *[]reflect.Value, method reflect.Method, jsonparams string) {
	v := pusher.New(reflectParams)
	params := createParameterArray(method)

	e := json.Unmarshal([]byte(jsonparams), &params)
	if e != nil {
		panic(e)
	}

	for _, param := range params {
		v.Push(reflect.ValueOf(param))
	}
}

func createUnmarshallingInterface(pointerType reflect.Type) interface{} {
	ok := pointerType.Kind() == reflect.Ptr
	if !ok {
		panic("Parameters have to be pointers")
	}

	valueType := pointerType.Elem()
	valuePointer := reflect.New(valueType)
	return valuePointer.Interface()
}

func executeCall(rcv reflect.Value, method reflect.Method, jsonparams string) []interface{} {
	params := []reflect.Value{rcv}
	decodeParams(&params, method, jsonparams)
	reflectresults := method.Func.Call(params)
	result := value2Interface(reflectresults)
	return result
}

func interface2Value(in []interface{}) []reflect.Value {
	out := make([]reflect.Value, 0)
	out_v := pusher.New(&out)
	for _, i := range in {
		out_v.Push(reflect.ValueOf(i))
	}
	return out
}

func value2Interface(in []reflect.Value) []interface{} {
	out := make([]interface{}, 0)
	out_v := pusher.New(&out)
	for _, i := range in {
		out_v.Push(i.Interface())
	}
	return out
}
