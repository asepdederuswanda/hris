// Package docs menyediakan API documentation menggunakan Scalar UI.
//
// Scalar adalah API documentation UI modern yang ringan dan cepat.
// Digunakan untuk mendokumentasikan endpoint REST API HRIS Platform.
package docs

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OpenAPISpec berisi spesifikasi OpenAPI 3.0 untuk HRIS Platform.
//
//go:embed openapi.json
var openapiFS embed.FS

// OpenAPIHandler mengembalikan Gin handler yang menyajikan file openapi.json.
func OpenAPIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := openapiFS.ReadFile("openapi.json")
		if err != nil {
			c.String(http.StatusNotFound, "OpenAPI spec not found")
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	}
}

// ScalarHTML adalah template HTML untuk Scalar UI.
const ScalarHTML = `<!doctype html>
<html>
<head>
    <title>HRIS Platform API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <div id="scalar-app" class="scalar-app"></div>
    <script
        id="api-reference"
        data-url="/openapi.json"
        data-hide-download-button="true"
        data-dark-mode="true"
        data-show-sidebar="true"
        data-force-dark-mode="true"
        type="application/json"
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

// ScalarUIHandler mengembalikan Gin handler yang menyajikan Scalar UI.
func ScalarUIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, ScalarHTML)
	}
}
