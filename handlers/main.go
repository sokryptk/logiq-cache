package handlers

import (
	"encoding/json"
	"logiq.ai/cache/cache"
	"logiq.ai/cache/models"
	"net/http"
)

func Store(c *cache.Cache, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	defer r.Body.Close()
	var entry models.Entry

	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		write(w, 401, "invalid request")
		return
	}

	c.Store(entry.Key, entry.Value)
	write(w, 200, "success")
}

func Retrieve(c *cache.Cache, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var entry models.Entry

	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		write(w, 401, "invalid request")
		return
	}

	value, err := c.Retrieve(entry.Key)
	if err != nil {
		write(w, 403, "invalid key")
		return
	}

	result := models.Response{
		Result: value,
	}

	write(w, 200, result)
}

func Delete(c *cache.Cache, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var entry models.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		write(w, 401, "invalid request")
		return
	}

	ok := c.Delete(entry.Key)
	if !ok {
		write(w, 403, "no such key")
		return
	}

	write(w, 200, "success")
}

func write(w http.ResponseWriter, code int, message any) {
	response, err := json.Marshal(message)
	if err != nil {
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}