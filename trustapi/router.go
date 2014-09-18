package trustapi

import "github.com/gorilla/mux"

func NewRouter(host string) *mux.Router {
	r := mux.NewRouter()
	v1 := r.Host(host).PathPrefix("/api/v1").Subrouter()
	v1.Path("/graph/base/{graphname}").Name("graphbase")
	return r
}
