package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func New() *httprouter.Router {

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Service is running"))
	})



	return router
}
