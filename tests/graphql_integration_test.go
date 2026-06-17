package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGraphQLIntegration(t *testing.T) {
	if _, err := os.Stat("../velonetics"); err != nil {
		t.Skip("velonetics binary not built; run make build first")
	}

	runner, _, err := NewIntegration(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer runner.Close()

	t.Run("mutationPOSTAdapter", func(t *testing.T) {
		body := `{"review":{"stars":5,"commentary":"great"}}`
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/graphql_adapter/review/1500", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := runner.httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		rawBody, ok := payload["body"].(string)
		if !ok {
			t.Fatalf("missing echoed body: %v", payload)
		}
		if !strings.Contains(rawBody, "CreateReviewForEpisode") {
			t.Fatalf("expected GraphQL mutation in upstream body: %s", rawBody)
		}
		if !strings.Contains(rawBody, `"stars":5`) {
			t.Fatalf("expected user review stars in upstream body: %s", rawBody)
		}
	})

	t.Run("queryGETAdapter", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/graphql_adapter/hero/JEDI", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := runner.httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		body := string(raw)
		if !strings.Contains(body, "Hero") {
			t.Fatalf("expected Hero query in upstream echo: %s", body)
		}
		if !strings.Contains(body, "JEDI") {
			t.Fatalf("expected JEDI variable in upstream echo: %s", body)
		}
	})

	t.Run("proxyPassthrough", func(t *testing.T) {
		clientBody := `{"query":"{ ping }","variables":{"a":1},"operationName":"Ping"}`
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/graphql_proxy", strings.NewReader(clientBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := runner.httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		rawBody, ok := payload["body"].(string)
		if !ok || rawBody != clientBody {
			t.Fatalf("proxy altered body: got %q want %q", rawBody, clientBody)
		}
	})

	t.Run("federationMerge", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/graphql_federation/42", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := runner.httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
		var payload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if _, ok := payload["user"]; !ok {
			t.Fatalf("missing user group: %v", payload)
		}
		if _, ok := payload["user_metadata"]; !ok {
			t.Fatalf("missing user_metadata group: %v", payload)
		}
	})

	t.Run("jwtBlocksUnauthorized", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/graphql_private/review/1", bytes.NewReader([]byte("{}")))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := runner.httpClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("expected auth failure, got 200")
		}
	})

	t.Run("rateLimitBlocksBurst", func(t *testing.T) {
		do := func() int {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/graphql_ratelimit", bytes.NewReader([]byte("{}")))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := runner.httpClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body)
			return resp.StatusCode
		}
		first := do()
		second := do()
		if first != http.StatusOK {
			t.Fatalf("first request status %d", first)
		}
		if second == http.StatusOK {
			time.Sleep(1100 * time.Millisecond)
			third := do()
			if third != http.StatusOK {
				t.Fatalf("expected recovery after window, got %d", third)
			}
			return
		}
		if second != http.StatusTooManyRequests && second != http.StatusInternalServerError {
			t.Fatalf("unexpected rate limit status %d", second)
		}
	})
}
