package portalgateway

import (
	"context"
	"fmt"
	"net/http"
	"portal-gateway/api"
	"portal-gateway/config"
	"portal-gateway/log"
	"portal-gateway/service"

	"github.com/julienschmidt/httprouter"
)

type Portal struct {
	service *httprouter.Router
	conf    *config.Config
	log     log.Logger
}

func NewPortal(conf *config.Config, log log.Logger) (*Portal, error) {
	ctx := context.Background()

	serviceRegistry, err := service.NewServiceRegistry(ctx, conf.GlobalConfig.ServiceType, conf)
	if err != nil {
		return nil, err
	}

	servicesRouter := api.NewServiceRouter(serviceRegistry)

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
