package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"logiq.ai/cache/cache"
)

type Entry struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Response struct {
	Result any
}

var c *cache.Cache

func main() {
	c = cache.New(1000)
	c.Logging = true
	c.SetTTL(time.Second * 5)
	mux := http.NewServeMux()
	mux.HandleFunc("/store", Store)
	mux.HandleFunc("/retrieve", Retrieve)
	mux.HandleFunc("/delete", Delete)

	log.Print("Running HTTP server at :4444")
	if err := http.ListenAndServe(":4444", mux); err != nil {
		log.Fatal(err)
	}
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	defer r.Body.Close()
	var entry Entry

	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		w.WriteHeader(403)
		w.Write([]byte("Invalid Request"))
		return
	}

	c.Store(entry.Key, entry.Value)
	w.Write([]byte("success"))
}

func Retrieve(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var entry Entry

	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("Invalid Request"))
		return
	}

	value, err := c.Retrieve(entry.Key)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("invalid key"))
		return
	}

	result := Response{
		Result: value,
	}

	// We can safeuly ignore this error since we exactly know what data
	// is being passed.
	resByte, _ := json.Marshal(result)

	w.Write(resByte)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		w.WriteHeader(405)
		w.Write([]byte("Invalid Request"))
		return
	}

	ok := c.Delete(entry.Key)
	if !ok {
		w.WriteHeader(405)
		w.Write([]byte("no such key"))
		return
	}

	w.Write([]byte("success"))
}