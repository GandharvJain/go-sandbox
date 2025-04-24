package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"fmt"
	"html/template"
)

const pagesRoot = "Pages/"
const pagesExt = ".txt"

type Page struct {
	Title	string
	Body	[]byte
}

func (p *Page) save() error {
	filename := pagesRoot + p.Title + pagesExt
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := pagesRoot + title + pagesExt
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/echo/*any", echoHandler)
	router.GET("/view/:title", makeHandler(viewHandler))
	router.GET("/edit/:title", makeHandler(editHandler))
	router.POST("/save/:title", makeHandler(saveHandler))
	router.Run("localhost:8080")
}

func makeHandler(fn func(*gin.Context, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		fn(c, title)
	}
}

func prerenderViewHandler(body []byte) []byte {
	re := regexp.MustCompile(`\[(.+)\]`)
	escapedBody := []byte(template.HTMLEscapeString(string(body)))
	newBody := re.ReplaceAllFunc(escapedBody, func(s []byte) []byte {
		pageName := re.ReplaceAllString(string(s), `$1`)
		newStr := fmt.Sprintf("<a href='/view/%s'>%s</a>", pageName, pageName)
		return []byte(newStr)
	})
	return newBody
}

func viewHandler(c *gin.Context, title string) {
	p, err := loadPage(title)
	if err != nil {
		c.HTML(http.StatusOK, "view_not_found.html", gin.H{"Title": title})
		return
	}
	body := string(prerenderViewHandler(p.Body))
	c.HTML(http.StatusOK, "view_found.html", gin.H{"Title": p.Title, "Body": template.HTML(body)})
}

func editHandler(c *gin.Context, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	c.HTML(http.StatusOK, "edit.html", p)
}

func saveHandler(c *gin.Context, title string) {
	body := c.PostForm("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusFound, "/view/"+title)
}

func echoHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hi there, I love %s!", c.Param("any")[1:])
}
