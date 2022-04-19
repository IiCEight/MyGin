package main

import (
	"fmt"
	"log"
	"net/http"
	"saber"
	"time"
)

func onlyForV2() saber.HandlerFunc {
	return func(c *saber.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		fmt.Println("mid")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := saber.Default()
	r.Static("/static", "D:/Project/Visual Studio Code/web/Mywebsite/about_me/style")

	alter := r.Group("/alter")
	// alter.Use(onlyForV2())
	alter.GET("/:name", func(c *saber.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	// r.SetFuncMap(template.FuncMap{"FormatAsDate": FormatAsDate})
	r.LoadHTMLGlob("D:/Project/Visual Studio Code/web/Mywebsite/about_me/*")
	r.GET("index", func(c *saber.Context) {
		c.HTML(http.StatusOK, "about_me.html", nil)
	})

	r.GET("/panic", func(c *saber.Context) {
		names := []string{"saber"}
		fmt.Println(c.Path)
		c.String(http.StatusOK, names[100])
	})

	r.Run("127.0.0.1:8080")
}
