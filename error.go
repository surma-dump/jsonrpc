package jsrpc

import (
	"web"
	"os"
	paths "path"
)

func serveError(ctx *web.Context, code int) {
	ctx.ContentType("html")
	serveTemplatedFile(ctx, paths.Join(prefix, "/errors/404.html"))
}

func getCodeByError(e os.Error) int {
	switch e {
	case os.EPERM:
		return 403
	default:
		return 404
	}
	return 0
}
