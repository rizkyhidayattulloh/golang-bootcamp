package http

import (
	"bytes"
	"kasir-api/docs"
	"net/http"
	"time"
)

var swaggerJSON = mustReadDocsFile("swagger.json")
var swaggerUIHTML = mustReadDocsFile("swagger-ui.html")

func mustReadDocsFile(name string) []byte {
	data, err := docs.EmbedAssets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}

type DocsController struct{}

func NewDocsController() *DocsController {
	return &DocsController{}
}

func (dc *DocsController) SwaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/swagger/doc.json" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	http.ServeContent(w, r, "swagger.json", time.Time{}, bytes.NewReader(swaggerJSON))
}

func (dc *DocsController) SwaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "swagger-ui.html", time.Time{}, bytes.NewReader(swaggerUIHTML))
}
