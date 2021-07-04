package main

import (
	_ "embed"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"regexp"
)

func Serve(dir string, port uint) {
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, &fileServerHandler{root: dir})
}

type fileServerHandler struct {
	root string
}

func (h *fileServerHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Note: not secure and prone to path manipulation attacks. Since it's a
	// development server, that's not particularly important, but it's worth
	// examining in the future.
	file := path.Join(h.root, req.URL.Path)
	ext := path.Ext(file)

	b, err := os.ReadFile(file)
	if os.IsNotExist(err) {
		http.NotFound(res, req)
		return
	} else if err != nil {
		log.Fatal(err)
	}

	if ext == ".html" {
		b = injectClientScripts(b)
	}

	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Content-Type", mime.TypeByExtension(ext))
	res.Header().Set("Content-Length", fmt.Sprint(len(b)))
	res.Write(b)
}

//go:embed client.js
var client []byte

func injectClientScripts(b []byte) []byte {
	bodyClosingTag := regexp.MustCompile("</body>")
	loc := bodyClosingTag.FindIndex(b)
	var pos int

	if loc != nil {
		pos = loc[0]
	} else {
		pos = len(b)
	}

	return insert(b, client, pos)
}

func insert(full, part []byte, pos int) []byte {
	inserted := append(full[:pos], part...)
	inserted = append(inserted, full[pos:]...)

	return inserted
}

func init() {
	client = append([]byte("<script type=\"module\">"), client...)
	client = append(client, []byte("</script>")...)
}
