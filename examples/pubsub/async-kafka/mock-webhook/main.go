package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	mu      sync.Mutex
	lastMsg map[string]interface{}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		lastMsg = body
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/last", func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if lastMsg == nil {
			http.Error(w, "no message yet", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(lastMsg)
	})
	http.ListenAndServe(":8081", mux)
}
