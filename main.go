package main

import (
	"frame"
	"net/http"
)

func main() {
	r := frame.New()
	r.GET("/", func(c *frame.Context) {
		c.String(http.StatusOK, "Hello qx32871\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *frame.Context) {
		names := []string{"qx32871"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":9999")
}
