package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	path    string
	fun     http.HandlerFunc
	methods []string
}

//	정적 URI Path
const StaticPath = "/static/"

//	로컬 경로
const StaticLocalDirPath = "./static/"

func makeHandlers(hs []handler, r *mux.Router, prefix string) {

	r.PathPrefix(StaticPath).Handler(http.StripPrefix(StaticPath, http.FileServer(http.Dir(StaticLocalDirPath))))
	apiRouter := r.PathPrefix(prefix).Subrouter()

	for _, h := range hs {
		if len(h.methods) == 0 {
			apiRouter.PathPrefix(h.path).HandlerFunc(h.fun)
		} else {
			apiRouter.PathPrefix(h.path).HandlerFunc(h.fun).Methods(h.methods...)
		}
	}

}
