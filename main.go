package main

import (
	"log"
	"net/http"
	"time"

	"logiq.ai/cache/cache"
	"logiq.ai/cache/handlers"
)

func main() {
	LRU := cache.New(1000)
	LRU.Logging = true
	LRU.SetExpiration(time.Second * 5)

	mux := http.NewServeMux()
	mux.HandleFunc("/store", supplyCache(LRU, handlers.Store))
	mux.HandleFunc("/retrieve", supplyCache(LRU, handlers.Retrieve))
	mux.HandleFunc("/delete", supplyCache(LRU, handlers.Delete))

	log.Print("Running HTTP server at :4444")
	if err := http.ListenAndServe(":4444", mux); err != nil {
		log.Fatal(err)
	}
}

// Inject Cache into our handlers. (Dependency Injection)
// This allows us to be more flexible with our cache in the future
// Also allows for mock testing.
func supplyCache(c *cache.Cache, handler func(c *cache.Cache, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(c, w, r)
	})
}