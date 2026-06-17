package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		resp, err := http.Get("http://127.0.0.1:4000/health")
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		return
	}

	addr := envOr("LISTEN_ADDR", ":4000")
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/graphql", graphqlHandler)
	mux.HandleFunc("/graphql/user", userSubgraphHandler)
	mux.HandleFunc("/graphql/metadata", metadataSubgraphHandler)

	log.Printf("graphql mock backend listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type graphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func readGraphQLRequest(r *http.Request) (graphQLRequest, []byte, error) {
	var req graphQLRequest
	if r.Method == http.MethodGet {
		req.Query = r.URL.Query().Get("query")
		req.OperationName = r.URL.Query().Get("operationName")
		if vars := r.URL.Query().Get("variables"); vars != "" {
			_ = json.Unmarshal([]byte(vars), &req.Variables)
		}
		raw, _ := json.Marshal(req)
		return req, raw, nil
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return req, nil, err
	}
	_ = json.Unmarshal(body, &req)
	return req, body, nil
}

func writeGraphQLResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	req, raw, err := readGraphQLRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("graphql proxy request method=%s operation=%s body=%s", r.Method, req.OperationName, string(raw))

	if strings.Contains(req.Query, "createReview") {
		writeGraphQLResponse(w, map[string]interface{}{
			"createReview": map[string]interface{}{
				"stars":      5,
				"commentary": "ok",
			},
		})
		return
	}
	if strings.Contains(req.Query, "Hero") {
		writeGraphQLResponse(w, map[string]interface{}{
			"hero": map[string]interface{}{"name": "Luke"},
		})
		return
	}
	writeGraphQLResponse(w, map[string]interface{}{"ping": "pong"})
}

func userSubgraphHandler(w http.ResponseWriter, r *http.Request) {
	req, raw, err := readGraphQLRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("user subgraph method=%s operation=%s body=%s", r.Method, req.OperationName, string(raw))
	writeGraphQLResponse(w, map[string]interface{}{
		"user": map[string]interface{}{"name": "Alice"},
	})
}

func metadataSubgraphHandler(w http.ResponseWriter, r *http.Request) {
	req, raw, err := readGraphQLRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("metadata subgraph method=%s operation=%s body=%s", r.Method, req.OperationName, string(raw))
	writeGraphQLResponse(w, map[string]interface{}{
		"metadata": map[string]interface{}{"locale": "en-US"},
	})
}
