package middlewares

import (
	"compress/gzip"
	"compress/zlib"

	"github.com/gin-gonic/gin"
)

// HTTPGzipEncoding ...
func HTTPGzipEncoding(c *gin.Context) {
	switch c.GetHeader("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			panic(err)
		}

		c.Request.Body = reader

	case "deflate":
		reader, err := zlib.NewReader(c.Request.Body)
		if err != nil {
			panic(err)
		}

		c.Request.Body = reader
	}

	c.Next()
}
