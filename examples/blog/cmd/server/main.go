package main

import (
	"cmp"
	"log"
	"net/http"
	"os"

	"github.com/typelate/dom/examples/blog/internal/hypertext"
)

func main() {
	mux := http.NewServeMux()
	app := new(hypertext.App)
	hypertext.TemplateRoutes(mux, app)
	log.Fatal(http.ListenAndServe(":"+cmp.Or(os.Getenv("PORT"), "8080"), mux))
}
