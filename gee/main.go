package main

import (
	"fmt"
	"net/http"

	"github.com/Muskchen/toolkits/gee/gee"
)

func main() {
	r := gee.Default()
	r.Use(gee.Logger())
	r.Get("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>", gee.H{})
	})

	v1 := r.Group("/v1")
	v1.Get("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee V1</h1>\"", gee.H{})
	})
	v1.Get("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	v1.Get("/hello/*filepath", func(c *gee.Context) {
		c.String(http.StatusOK, c.Param("filepath"))
	})

	v2 := r.Group("/v2")
	v2.Get("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee V2</h1>\"", gee.H{})
	})
	v2.Get("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	if err := r.Run(":9999"); err != nil {
		fmt.Println(err)
	}
}
