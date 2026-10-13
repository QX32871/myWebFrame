package main

import (
	"frame"
	"log"
	"net/http"
	"time"
)

func onlyForV2() frame.HandlerFunc {
	return func(c *frame.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := frame.New()
	r.Use(frame.Logger())
	r.GET("/", func(c *frame.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *frame.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	err := r.Run(":9999")
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
		return
	}
}
