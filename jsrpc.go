package jsrpc

import (
	"web"
	paths "path"
	"reflect"
)

type callback struct {
	object  reflect.Value
	methods map[string]reflect.Method
}

var (
	cbos   map[string]callback = make(map[string]callback)
	prefix string
)

// Starts the webserver on given port, serving the file
// in given path. If standalone is false, incoming connections
// are expected to be FCGI, raw HTTP otherwise.
func Run(path string, addr string, standalone bool) {
	prefix = path
	web.Get("(.*)", handler)
	if standalone {
		web.Run(addr)
	} else {
		web.RunFcgi(addr)
	}
}

// Registers an objects public methods for JSRPC.
// Every path beginning with given prefix will be
// parsed to a function call which has to match the
// on of the public functions defined on v.
// Old values can be overwritten.
func RegisterRPC(prefix string, v interface{}) {
	cbos[paths.Clean(prefix)] = parseInterface(v)
}

func handler(ctx *web.Context, path string) {
	cboname := paths.Clean(dirname(path))

	cbo, ok := cbos[cboname]
	if ok {
		methodname := paths.Base(path)
		result := callObject(cbo, methodname, ctx.Params["params"])
		ctx.WriteString(jsonify(result))
	} else {
		serveFile(ctx, path)
	}
}

func callObject(cbo callback, methodname string, params string) (result interface{}) {
	// Catch the special function "_enumerate" which lists all 
	// available methods
	if methodname == "_enumerate" {
		result = generateMethodEnumeration(cbo)
	} else {
		cbf, ok := cbo.methods[methodname]
		if !ok {
			// FIXME Does this need more
			// gentle handling?
			panic("Unknown function")
		}
		result = executeCall(cbo.object, cbf, params)
	}
	return
}

type methodDescription struct {
	Name  string
	NumIn int
}

func generateMethodEnumeration(cbo callback) (m []methodDescription) {
	m = make([]methodDescription, len(cbo.methods))
	i := 0
	for name, method := range cbo.methods {
		m[i].Name = name
		_ = method
		// first parameter is the receiver object
		// which is provided by us, not the remote caller
		m[i].NumIn = method.Type.NumIn() - 1
		i++
	}
	return
}

func serveFile(ctx *web.Context, path string) {
	if path == "/" {
		path = "/index.html"
	}

	ext := paths.Ext(path)
	ctx.ContentType(ext)
	path = cleanPath(path)
	if isTemplateType(ext) {
		serveTemplatedFile(ctx, path)
	} else {
		serveStaticFile(ctx, path)
	}
}

func dirname(path string) string {
	dir, _ := paths.Split(path)
	return dir
}
