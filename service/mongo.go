package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type mongoServiceRegistry struct {
	*baseRegistry
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func NewMongoServiceRegistry(ctx context.Context, client *mongo.Client, database string, collection string, br *baseRegistry) (ServiceRegistry, error) {
	r := &mongoServiceRegistry{
		baseRegistry: br,
		client:       client,
		database:     client.Database(database),
		ctx:          ctx,
	}
	r.collection = r.database.Collection(collection)
	err := r.loadServices()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *mongoServiceRegistry) AddService(service *BackendService) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.addService(service, func() error {
		_, err := r.collection.InsertOne(r.ctx, service)
		return err
	})
}

func (r *mongoServiceRegistry) GetServices() []*BackendService {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.getServices()
}

func (r *mongoServiceRegistry) GetService(name string) (*BackendService, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.getService(name)
}

func (r *mongoServiceRegistry) loadServices() error {
	cursor, err := r.collection.Find(r.ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(r.ctx)

	for cursor.Next(r.ctx) {
		var service BackendService
		err = cursor.Decode(&service)
		if err != nil {
			return err
		}

		r.services = append(r.services, &service)
	}

	if err = cursor.Err(); err != nil {
		return err
	}

	return nil
}
