package main

import (
	"fmt"
	"frame"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := frame.New()
	r.Use(frame.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("static/templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "QX32871", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *frame.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *frame.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", frame.H{
			"title":  "QX36871",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *frame.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", frame.H{
			"title": "QX34871",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
