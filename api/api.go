package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"portal-gateway/service"

	"github.com/julienschmidt/httprouter"
)

func NewServiceRouter(backendServices service.ServiceRegistry) *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/ping", getHealth())

	router.GET("/api/services", getServicesHandler(backendServices))
	router.POST("/api/services", addServiceHandler(backendServices))
	router.GET("/api/service/:name", getServiceHandler(backendServices))
	return router
}

func getHealth() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		prepareHeaders(w, http.StatusOK)
		json.NewEncoder(w).Encode("pong!")
	}
}

func getServicesHandler(bs service.ServiceRegistry) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		services := bs.GetServices()
		jsonData, err := json.Marshal(services)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		prepareHeaders(w, http.StatusOK)
		w.Write(jsonData)
	}
}

func getServiceHandler(bs service.ServiceRegistry) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		service, err := bs.GetService(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData, err := json.Marshal(service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		prepareHeaders(w, http.StatusOK)
		w.Write(jsonData)

	}
}

func addServiceHandler(bs service.ServiceRegistry) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var service service.BackendService
		err := json.NewDecoder(r.Body).Decode(&service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = validateService(&service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = bs.AddService(&service)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		prepareHeaders(w, http.StatusCreated)
		json.NewEncoder(w).Encode(service)
	}
}

func validateService(service *service.BackendService) error {
	if service.Path == "" {
		return fmt.Errorf("path is a required field")
	}
	if len(service.UpstreamTargets) < 1 {
		return fmt.Errorf("at least one upstream target is required")
	}

	for _, target := range service.UpstreamTargets {
		u, err := url.Parse(target)
		if err != nil {
			return fmt.Errorf("Invalid upstream: " + target)
		}
		if u.Scheme == "" {
			return fmt.Errorf("Upstream " + target + " should be include a scheme (e.g., 'http' or 'https')")
		}
	}

	if service.Scheme == "" {
		service.Scheme = "http"
	}

	if service.Timeout == 0 {
		service.Timeout = 10
	}

	return nil
}

func prepareHeaders(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
