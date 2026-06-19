package pucora

import (
	"fmt"

	apikeys "github.com/pucora/pucora-apikeys/v2"
	apikeysgin "github.com/pucora/pucora-apikeys/v2/router/gin"
	basicauthingin "github.com/pucora/pucora-basicauth/v2/router/gin"
	botdetector "github.com/pucora/pucora-botdetector/v2/gin"
	jose "github.com/pucora/pucora-jose/v2"
	ginjose "github.com/pucora/pucora-jose/v2/gin"
	lua "github.com/pucora/pucora-lua/v2/router/gin"
	metrics "github.com/pucora/pucora-metrics/v2/gin"
	opencensus "github.com/pucora/pucora-opencensus/v2/router/gin"
	ratelimit "github.com/pucora/pucora-ratelimit/v3/router/gin"
	wsgin "github.com/pucora/pucora-websocket/v2/router/gin"
	"github.com/pucora/lura/v2/config"
	"github.com/pucora/lura/v2/logging"
	"github.com/pucora/lura/v2/proxy"
	router "github.com/pucora/lura/v2/router/gin"
	"github.com/pucora/lura/v2/transport/http/server"

	"github.com/gin-gonic/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory, serviceCfg config.ServiceConfig) router.HandlerFactory {
	basicSvc := basicauthingin.NewService(serviceCfg, logger)
	apiKeysCfg, ok := apikeys.ParseServiceConfig(serviceCfg.ExtraConfig)
	var apiKeysRegistry *apikeys.Registry
	if ok {
		apiKeysRegistry = apikeys.NewRegistry(apiKeysCfg)
	}

	handlerFactory := router.CustomErrorEndpointHandler(logger, server.DefaultToHTTPError)
	handlerFactory = ratelimit.NewRateLimiterMw(logger, handlerFactory)
	handlerFactory = lua.HandlerFactory(logger, handlerFactory)
	handlerFactory = wsgin.HandlerFactory(logger, handlerFactory)
	handlerFactory = apikeysgin.HandlerFactory(handlerFactory, logger, apiKeysRegistry)
	handlerFactory = basicauthingin.HandlerFactory(handlerFactory, logger, basicSvc)
	handlerFactory = ginjose.HandlerFactory(handlerFactory, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	handlerFactory = botdetector.New(handlerFactory, logger)

	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		logger.Debug(fmt.Sprintf("[ENDPOINT: %s] Building the http handler", cfg.Endpoint))
		return handlerFactory(cfg, p)
	}
}

type handlerFactory struct{}

func (handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory, sc config.ServiceConfig) router.HandlerFactory {
	return NewHandlerFactory(l, m, r, sc)
}
