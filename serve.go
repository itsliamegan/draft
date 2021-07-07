package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
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
	path := filepath.Join(h.root, req.URL.Path)
	stat, err := os.Stat(path)

	res.Header().Set("Cache-Control", "no-cache")

	if errors.Is(err, os.ErrNotExist) {
		http.NotFound(res, req)
	} else if err != nil {
		log.Fatal(err)
	} else if stat.IsDir() {
		h.serveDir(path, res, req)
	} else {
		h.serveFile(path, res, req)
	}
}

func (h *fileServerHandler) serveDir(dir string, res http.ResponseWriter, req *http.Request) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.Write([]byte("<table style=\"font-family: monospace;\"><tbody>"))

	for _, entry := range entries {
		relativePath, err := filepath.Rel(h.root, filepath.Join(dir, entry.Name()))

		if err != nil {
			log.Fatal(err)
		}

		var displayName string

		if entry.IsDir() {
			displayName = entry.Name() + "/"
		} else {
			displayName = entry.Name()
		}

		res.Write([]byte("<tr><td>"))
		anchor := fmt.Sprintf("<a href=\"%s\">%s</a>", relativePath, displayName)
		res.Write([]byte(anchor))
		res.Write([]byte("</td></tr>"))
	}
}

func (h *fileServerHandler) serveFile(file string, res http.ResponseWriter, req *http.Request) {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	ext := filepath.Ext(file)
	mimeType := mime.TypeByExtension(ext)

	if mimeType == "text/html" {
		b = injectClientScripts(b)
	} else if mimeType == "" {
		mimeType = "text/plain"
	}

	res.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", mimeType))
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
