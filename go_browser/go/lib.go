package main

/*
#include <stdlib.h>

// Define the C struct equivalent to the Go struct
typedef struct {
    char* Title;
    char* Link;
    char* Time;
    char* Author;
    char* SourceLink;
    char* SourceName;
} CArticle;
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/headzoo/surf.v1"
)

func main() {
	cArticles := scrapeGoogleNews()
	articlePrinter(cArticles)
}

func articlePrinter(articles *C.CArticle) {
	// Convert the Reference to a Pointer
	c := (*C.CArticle)(unsafe.Pointer(articles))
	fmt.Println(C.GoString(c.Title))
}

//export freeCArticles
func freeCArticles(cArticles *C.CArticle, size C.int) {
	c := (*C.CArticle)(unsafe.Pointer(cArticles))

	for i := 0; i < int(size); i++ {
		if c.Title != nil {
			C.free(unsafe.Pointer(c.Title))
			C.free(unsafe.Pointer(c.Link))
			C.free(unsafe.Pointer(c.Time))
			C.free(unsafe.Pointer(c.Author))
			C.free(unsafe.Pointer(c.SourceLink))
			C.free(unsafe.Pointer(c.SourceName))
		}

		// Move to the next CArticle
		c = (*C.CArticle)(unsafe.Pointer(uintptr(unsafe.Pointer(c)) + unsafe.Sizeof(*c)))
	}

	// Free the original pointer
	C.free(unsafe.Pointer(cArticles))
}

type Article struct {
	Title      string
	Link       string
	Time       string
	Author     string
	SourceLink string
	SourceName string
}

//export scrapeGoogleNews
func scrapeGoogleNews() *C.CArticle {
	bow := surf.NewBrowser()
	err := bow.Open("https://news.google.com/topics/CAAqJggKIiBDQkFTRWdvSUwyMHZNRFZxYUdjU0FtVnVHZ0pWVXlnQVAB?hl=en-US&gl=US&ceid=US%3Aen")
	if err != nil {
		panic(err)
	}
	max_count := 100
	articles := make([]Article, 0, max_count)
	bow.Dom().Find("article").Each(func(_ int, s *goquery.Selection) {
		title := s.Find("a").Text()
		time := s.Find("time").AttrOr("datetime", "")
		link := s.Find("a").AttrOr("href", "")
		source_link := s.Find("div[data-n-tid='9']").First().Parent().Parent().Find("img").First().AttrOr("src", "")
		source_name := s.Find("div[data-n-tid='9']").First().Text()
		author := s.Find("time").Parent().Find("span").First().Text()
		if author == "" {
			author = "Syndicated Source"
		}
		article := Article{
			Title:      title,
			Link:       link,
			Time:       time,
			Author:     author,
			SourceLink: source_link,
			SourceName: source_name,
		}
		articles = append(articles, article)
	})
	fmt.Println(len(articles))
	// Convert Go articles to CArticles
	cArticlesArray := C.malloc(C.size_t(len(articles)) * C.size_t(unsafe.Sizeof(C.CArticle{})))
	cArraySlice := (*[1 << 30]C.CArticle)(cArticlesArray)[:len(articles):len(articles)]
	for i, article := range articles {
		cArraySlice[i] = C.CArticle{
			Title:      C.CString(article.Title),
			Link:       C.CString(article.Link),
			Time:       C.CString(article.Time),
			Author:     C.CString(article.Author),
			SourceLink: C.CString(article.SourceLink),
			SourceName: C.CString(article.SourceName),
		}
	}
	return (*C.CArticle)(cArticlesArray)
}

//export ExampleBrowser
func ExampleBrowser(path *C.char) *C.char {
	bow := surf.NewBrowser()
	err := bow.Open(C.GoString(path))
	if err != nil {
		panic(err)
	}
	fmt.Println(bow.Title())
	return C.CString(fmt.Sprintf("Hello from Go, %s!", C.GoString(path)))
}
