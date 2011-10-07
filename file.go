package jsrpc

import (
	"web"
	"mustache"
	paths "path"
	"io"
	"os"
	"time"
)

func isTemplateType(extension string) bool {
	return extension == ".html" || extension == ".htm"
}
func cleanPath(path string) string {
	path = paths.Clean(path)
	return paths.Clean(paths.Join(prefix, "/", path))
}

func serveTemplatedFile(ctx *web.Context, path string) {
	ctx.WriteString(mustache.RenderFile(path, getTemplateEnvironment()))
}

func getTemplateEnvironment() map[string]string {
	return map[string]string{
		"Generationdate": time.UTC().String(),
	}
}

func serveStaticFile(ctx *web.Context, path string) {
	f, e := os.Open(path)
	if e != nil {
		serveError(ctx, getCodeByError(e))
		return
	}
	defer f.Close()

	io.Copy(ctx, f)
}
