package service

import (
	"context"
	"portal-gateway/config"
	"sync"
)

type ServiceRegistry interface {
	GetServices() []*BackendService
	AddService(service *BackendService) error
	GetService(name string) (*BackendService, error)
	UpdateService(name string, service *BackendService) error
	RemoveService(name string) error
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

func (r *baseRegistry) getService(name string) (*BackendService, error) {
	var srv *BackendService

	for i, s := range r.services {
		if s.Name == name {
			srv = r.services[i]
			return srv, nil
		}
	}
	return nil, ErrServiceNotFound{Name: name}
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

func (r *baseRegistry) updateService(name string, service *BackendService, apply func() error) error {
	old := r.getServices()

	for i, s := range r.services {
		if s.Name == name {
			r.services[i] = service
			err := apply()
			if err != nil {
				r.services = old
				return err
			}

			return nil
		}
	}

	return ErrServiceNotFound{Name: name}
}

func (r *baseRegistry) removeService(name string, apply func() error) error {
	old := r.getServices()

	for i, s := range r.services {
		if s.Name == name {
			r.services = append(r.services[:i], r.services[i+1:]...)
			err := apply()
			if err != nil {
				r.services = old
				return err
			}

			return nil
		}
	}

	return ErrServiceNotFound{Name: name}
}
