package web

import (
	"embed"
	"path"
)

//go:embed admin client/dist static
var FS embed.FS

var (
	CLIENT_JS  = mustAsset(".js")
	CLIENT_CSS = mustAsset(".css")
)

// the bundle carries a content hash in its name, the server looks it up once
func mustAsset(ext string) string {
	files, err := FS.ReadDir("client/dist")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if path.Ext(f.Name()) == ext {
			return f.Name()
		}
	}
	panic("client " + ext + " asset not found")
}
