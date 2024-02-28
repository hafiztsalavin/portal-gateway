package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewServiceRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/api/ping", getHealth())
	return router
}

func getHealth() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		prepareHeaders(w, http.StatusOK)
		json.NewEncoder(w).Encode("pong!")
	}

}

func prepareHeaders(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}
