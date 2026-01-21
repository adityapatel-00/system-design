package main

import (
	"net/http"

	"github.com/adityapatel-00/system-design/design-problems/easy/go/urlshortner/handler"
	"github.com/adityapatel-00/system-design/design-problems/easy/go/urlshortner/store"
)

func main() {
	server := http.NewServeMux()

	store := store.NewURLStore()

	server.HandleFunc("/shorten", handler.SaveNewURL(store))
	server.HandleFunc("/analytics/", handler.GetAnalytics(store))
	server.HandleFunc("/", handler.RedirectUrl(store))

	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(err)
	}
}
