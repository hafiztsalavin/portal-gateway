package portalgateway

import (
	"fmt"
	"net/http"
	"portal-gateway/api"
	"portal-gateway/config"
	"portal-gateway/log"

	"github.com/julienschmidt/httprouter"
)

type Portal struct {
	service *httprouter.Router
	conf    *config.Config
	log     log.Logger
}

func NewPortal(conf *config.Config, log log.Logger) (*Portal, error) {
	servicesRouter := api.NewServiceRouter()

	return &Portal{
		service: servicesRouter,
		conf:    conf,
		log:     log,
	}, nil
}

func (gw *Portal) Start() error {
	gatewayAddr := gw.conf.GatewayConfig.Addr
	if gatewayAddr == "" {
		gatewayAddr = "0.0.0.0:8000"
	}

	gatewayHandler := gw.service
	gw.log.WithFields(log.InfoLevel, fmt.Sprintf("Started on %s", gatewayAddr))
	if err := startServer(gatewayAddr, gatewayHandler); err != nil {
		return err
	}

	return nil
}

func startServer(addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
