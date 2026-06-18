package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestPubSubMemRoundTripIntegration(t *testing.T) {
	bin := "../velonetics"
	if _, err := os.Stat(bin); err != nil {
		t.Skip("velonetics binary not built; run make build first")
	}

	gatewayPort, err := freePort()
	if err != nil {
		t.Fatal(err)
	}

	cfg := map[string]interface{}{
		"version": 3,
		"port":    gatewayPort,
		"extra_config": map[string]interface{}{
			"telemetry/logging": map[string]interface{}{
				"level":  "ERROR",
				"stdout": true,
			},
			"telemetry/usage": map[string]interface{}{
				"enabled": false,
			},
		},
		"endpoints": []map[string]interface{}{
			{
				"endpoint": "/publish",
				"method":   "POST",
				"backend": []map[string]interface{}{
					{
						"host":                  []string{"mem://smoke-host"},
						"url_pattern":           "/ignored",
						"disable_host_sanitize": true,
						"extra_config": map[string]interface{}{
							"backend/pubsub/publisher": map[string]interface{}{
								"topic_url": "/events",
							},
						},
					},
				},
			},
			{
				"endpoint": "/subscribe",
				"method":   "GET",
				"backend": []map[string]interface{}{
					{
						"host":                  []string{"mem://smoke-host"},
						"url_pattern":           "/ignored",
						"disable_host_sanitize": true,
						"encoding":              "json",
						"extra_config": map[string]interface{}{
							"backend/pubsub/subscriber": map[string]interface{}{
								"subscription_url": "/events",
							},
						},
					},
				},
			},
		},
	}

	cfgPath := filepath.Join(t.TempDir(), "pubsub_integration.json")
	f, err := os.Create(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cmd := exec.Command(bin, "run", "-c", cfgPath)
	cmd.Env = append(os.Environ(), "USAGE_DISABLE=1")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()

	healthURL := fmt.Sprintf("http://127.0.0.1:%d/__health", gatewayPort)
	if err := waitForReady(healthURL, 15*time.Second); err != nil {
		t.Fatalf("gateway not ready: %v", err)
	}

	base := fmt.Sprintf("http://127.0.0.1:%d", gatewayPort)
	pubReq, err := http.NewRequest(http.MethodPost, base+"/publish", bytes.NewBufferString(`{"event":"integration-test"}`))
	if err != nil {
		t.Fatal(err)
	}
	pubReq.Header.Set("Content-Type", "application/json")
	pubResp, err := http.DefaultClient.Do(pubReq)
	if err != nil {
		t.Fatal(err)
	}
	pubResp.Body.Close()
	if pubResp.StatusCode != http.StatusOK {
		t.Fatalf("publish status: %d", pubResp.StatusCode)
	}

	var subResp *http.Response
	for i := 0; i < 20; i++ {
		subResp, err = http.Get(base + "/subscribe")
		if err == nil && subResp.StatusCode == http.StatusOK {
			break
		}
		if subResp != nil {
			subResp.Body.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer subResp.Body.Close()

	body, err := io.ReadAll(subResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		t.Fatalf("invalid json: %s", body)
	}
	if data["event"] != "integration-test" {
		t.Fatalf("unexpected payload: %v", data)
	}
}

func TestPubSubFixtureSchema(t *testing.T) {
	bin := "../velonetics"
	if _, err := os.Stat(bin); err != nil {
		t.Skip("velonetics binary not built; run make build first")
	}
	path := filepath.Join("fixtures", "pubsub_mem.json")
	cmd := exec.Command(bin, "check", "-c", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check failed: %v\n%s", err, out)
	}
}

func TestPubSubFixtureUsesMemDriver(t *testing.T) {
	path := filepath.Join("fixtures", "pubsub_mem.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`mem://`)) {
		t.Fatalf("%s: expected mem:// backend host", path)
	}
	if !bytes.Contains(data, []byte(`backend/pubsub/publisher`)) {
		t.Fatalf("%s: missing publisher extra_config", path)
	}
	if !bytes.Contains(data, []byte(`backend/pubsub/subscriber`)) {
		t.Fatalf("%s: missing subscriber extra_config", path)
	}
}
