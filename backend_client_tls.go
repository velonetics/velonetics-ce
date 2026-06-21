package pucora

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pucora/lura/v2/config"
	"github.com/pucora/lura/v2/logging"
	"github.com/pucora/lura/v2/transport/http/client"
	httpserver "github.com/pucora/lura/v2/transport/http/server"
)

const backendHTTPClientNamespace = "backend/http/client"

func parseBackendHTTPClientTLS(cfg *config.Backend) (*config.ClientTLS, bool) {
	if cfg == nil || cfg.ExtraConfig == nil {
		return nil, false
	}
	v, ok := cfg.ExtraConfig[backendHTTPClientNamespace].(map[string]interface{})
	if !ok {
		return nil, false
	}
	tlsRaw, ok := v["client_tls"]
	if !ok {
		return nil, false
	}
	raw, err := json.Marshal(tlsRaw)
	if err != nil {
		return nil, false
	}
	var clientTLS config.ClientTLS
	if err := json.Unmarshal(raw, &clientTLS); err != nil {
		return nil, false
	}
	return &clientTLS, true
}

func newHTTPClientWithBackendTLS(cfg *config.Backend, next client.HTTPClientFactory, logger logging.Logger) client.HTTPClientFactory {
	clientTLS, ok := parseBackendHTTPClientTLS(cfg)
	if !ok {
		return next
	}
	tlsConfig := httpserver.ParseClientTLSConfigWithLogger(clientTLS, logger)
	if tlsConfig == nil {
		return next
	}
	return func(ctx context.Context) *http.Client {
		base := next(ctx)
		transport := cloneHTTPTransport(base.Transport)
		transport.TLSClientConfig = tlsConfig
		return &http.Client{
			Transport:     transport,
			CheckRedirect: base.CheckRedirect,
			Jar:           base.Jar,
			Timeout:       base.Timeout,
		}
	}
}

func cloneHTTPTransport(rt http.RoundTripper) *http.Transport {
	if rt == nil {
		if dt, ok := http.DefaultTransport.(*http.Transport); ok {
			return dt.Clone()
		}
		return &http.Transport{}
	}
	if t, ok := rt.(*http.Transport); ok {
		return t.Clone()
	}
	if dt, ok := http.DefaultTransport.(*http.Transport); ok {
		return dt.Clone()
	}
	return &http.Transport{}
}
