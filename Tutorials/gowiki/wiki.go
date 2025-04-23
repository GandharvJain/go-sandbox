package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Page struct {
	Title	string
	Body	[]byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, err := loadPage("TestPage")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while loading page: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(p2.Body))

	router := gin.Default()
	router.GET("/echo/*any", echoBack)
	router.GET("/view/:title", viewPage)
	router.Run("localhost:8080")
}

func viewPage(c *gin.Context) {
	title := c.Param("title")
	p, err := loadPage(title)
	if err != nil {
		data := fmt.Sprintf("<h1>Page \"%s\" does not exist</h1>", title)
		c.Data(http.StatusNotFound, "text/html", []byte(data))
	}
	data := fmt.Sprintf("<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	c.Data(http.StatusOK, "text/html", []byte(data))
}

func echoBack(c *gin.Context) {
	c.String(http.StatusOK, "Hi there, I love %s!", c.Param("any")[1:])
}
