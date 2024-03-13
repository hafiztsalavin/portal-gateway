package service

import (
	"context"
	"portal-gateway/config"
	"sync"
)

type ServiceRegistry interface {
	GetServices() []*BackendService
	AddService(service *BackendService) error
}

type baseRegistry struct {
	mutex    *sync.RWMutex
	services []*BackendService
}

func NewServiceRegistry(ctx context.Context, serviceType string, config *config.Config) (ServiceRegistry, error) {
	var (
		reg ServiceRegistry
		mu  sync.RWMutex
		err error
	)
	baseReg := baseRegistry{
		mutex: &mu,
	}

	mongoClient, err := NewMongoClient(ctx, config.GlobalConfig.MongoURI)
	if err != nil {
		return nil, err
	}
	reg, err = NewMongoServiceRegistry(ctx, mongoClient, config.GlobalConfig.MongoDatabaseName, config.GlobalConfig.MongoCollectionName, &baseReg)
	if err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *baseRegistry) getServices() []*BackendService {
	services := make([]*BackendService, len(r.services))
	copy(services, r.services)
	return services
}

func (r *baseRegistry) addService(service *BackendService, apply func() error) error {
	old := r.getServices()

	for _, s := range r.services {
		if s.Name == service.Name {
			return ErrServiceExists{Name: service.Name}
		}
	}

	r.services = append(r.services, service)
	err := apply()
	if err != nil {
		r.services = old
		return err
	}

	return nil
}
