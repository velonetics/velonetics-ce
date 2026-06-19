// Pucora-ce sets up a complete Pucora API Gateway ready to serve

package main

import (
	"context"
	"embed"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pucora "github.com/pucora/velonetics-ce/v2"
	cmd "github.com/pucora/velonetics-cobra/v2"
	flexibleconfig "github.com/pucora/velonetics-flexibleconfig/v2"
	koanf "github.com/pucora/velonetics-koanf"
	"github.com/pucora/lura/v2/config"
)

const (
	fcPartials  = "FC_PARTIALS"
	fcTemplates = "FC_TEMPLATES"
	fcSettings  = "FC_SETTINGS"
	fcPath      = "FC_OUT"
	fcEnable    = "FC_ENABLE"
)

//go:embed schema
var embedSchema embed.FS

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	pucora.RegisterEncoders()

	registerAliases()

	var cfg config.Parser
	cfg = koanf.New()
	if os.Getenv(fcEnable) != "" {
		cfg = flexibleconfig.NewTemplateParser(flexibleconfig.Config{
			Parser:    cfg,
			Partials:  os.Getenv(fcPartials),
			Settings:  os.Getenv(fcSettings),
			Path:      os.Getenv(fcPath),
			Templates: os.Getenv(fcTemplates),
		})
	}

	var rawSchema string
	schema, err := embedSchema.ReadFile("schema/schema.json")
	if err == nil {
		rawSchema = string(schema)
	}

	commandsToLoad := []cmd.Command{
		cmd.RunCommand,
		cmd.NewCheckCmd(rawSchema),
		cmd.PluginCommand,
		cmd.VersionCommand,
		cmd.AuditCommand,
		pucora.NewTestPluginCmd(),
	}

	cmd.DefaultRoot = cmd.NewRoot(cmd.RootCommand, commandsToLoad...)
	cmd.DefaultRoot.Cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.Execute(cfg, pucora.NewExecutor(ctx))
}

var aliases = map[string]string{
	"github_com/pucora/velonetics-ce/v2/transport/http/server/handler":  "plugin/http-server",
	"github.com/pucora/velonetics-ce/v2/transport/http/client/executor": "plugin/http-client",
	"github.com/pucora/velonetics-ce/v2/proxy/plugin":                   "plugin/req-resp-modifier",
	"github.com/pucora/velonetics-ce/v2/proxy":                          "proxy",
	"github_com/pucora/lura/v2/router/gin":                              "router",

	"github.com/pucora/velonetics-httpcache":                "qos/http-cache",
	"github.com/pucora/velonetics-circuitbreaker/gobreaker": "qos/circuit-breaker",

	"github.com/pucora/velonetics-oauth2-clientcredentials": "auth/client-credentials",
	"github.com/pucora/velonetics-jose/validator":           "auth/validator",
	"github.com/pucora/velonetics-jose/signer":              "auth/signer",
	"github_com/devopsfaith/bloomfilter":                        "auth/revoker",

	"github_com/pucora/velonetics-botdetector": "security/bot-detector",
	"github_com/pucora/velonetics-httpsecure":  "security/http",
	"github_com/pucora/velonetics-cors":        "security/cors",

	"github.com/pucora/velonetics-cel":        "validation/cel",
	"github.com/pucora/velonetics-jsonschema": "validation/json-schema",

	"github.com/pucora/velonetics-amqp/agent": "async/amqp",

	"github.com/pucora/velonetics-amqp/consume":                  "backend/amqp/consumer",
	"github.com/pucora/velonetics-amqp/produce":                  "backend/amqp/producer",
	"github.com/pucora/velonetics-lambda":                        "backend/lambda",
	"github.com/pucora/velonetics-soap/v2":                       "backend/soap",
	"github.com/pucora/velonetics-grpc/v2":                       "grpc",
	"github.com/pucora/velonetics-grpc/v2/client":                "backend/grpc",
	"github.com/pucora/velonetics-pubsub/publisher":              "backend/pubsub/publisher",
	"github.com/pucora/velonetics-pubsub/subscriber":             "backend/pubsub/subscriber",
	"github.com/pucora/velonetics-pubsub/kafka/publisher":        "backend/pubsub/publisher/kafka",
	"github.com/pucora/velonetics-pubsub/kafka/subscriber":       "backend/pubsub/subscriber/kafka",
	"github.com/pucora/velonetics-pubsub/async":                "async/kafka",
	"github.com/pucora/lura/v2/transport/http/client/graphql": "backend/graphql",
	"github.com/pucora/velonetics-ce/v2/http":                          "backend/http",

	"github_com/pucora/velonetics-gelf":       "telemetry/gelf",
	"github_com/pucora/velonetics-gologging":  "telemetry/logging",
	"github_com/pucora/velonetics-logstash":   "telemetry/logstash",
	"github_com/pucora/velonetics-metrics":    "telemetry/metrics",
	"github_com/pucora/velonetics-influx":     "telemetry/influx",
	"github_com/pucora/velonetics-opencensus": "telemetry/opencensus",

	"github.com/pucora/velonetics-lua/router":        "modifier/lua-endpoint",
	"github.com/pucora/velonetics-lua/proxy":         "modifier/lua-proxy",
	"github.com/pucora/velonetics-lua/proxy/backend": "modifier/lua-backend",
	"github.com/pucora/velonetics-martian":           "modifier/martian",

}

func registerAliases() {
	for key, alias := range aliases {
		config.ExtraConfigAlias[alias] = key
	}
	// Legacy velonetics.io namespace keys in extra_config remain accepted.
	for key := range aliases {
		legacyKey := strings.Replace(key, "github.com/pucora/", "github.com/velonetics/", 1)
		legacyKey = strings.Replace(legacyKey, "github_com/pucora/", "github_com/velonetics/", 1)
		if legacyKey != key {
			config.ExtraConfigAlias[legacyKey] = key
		}
	}
}
