package main

import (
	"frame"
	"net/http"
)

func main() {
	router := frame.New()
	router.GET("/", func(c *frame.Context) {
		c.HTML(http.StatusOK, "<h1>Hello World!</h1>")
	})
	router.GET("/hello", func(c *frame.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	router.GET("/hello/:name", func(c *frame.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	router.POST("/assets/*filepath", func(c *frame.Context) {
		c.JSON(http.StatusOK, frame.H{
			"filepath": c.Param("filepath"),
		})
	})
	err := router.Run(":9999")
	if err != nil {
		return
	}
}
