package middlewares

import "net/http"

type middleware func(http.Handler) http.Handler

func ApplyMiddlewares(handler http.Handler, middlewares ...middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
