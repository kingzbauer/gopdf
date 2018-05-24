package server

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func RequirePOST(view http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		view.ServeHTTP(w, r)
	})
}

func WrapHandlerFunc(middleware Middleware, fun http.HandlerFunc) http.Handler {
	return middleware(fun)
}
