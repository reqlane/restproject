package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type HPPOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	WhiteList                   []string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckQuery && r.URL.Query() != nil {
				// Filter the query params
				filterQueryParams(r, options.WhiteList)
			}

			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				// Filter the body params
				filterBodyParams(r, options.WhiteList)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterQueryParams(r *http.Request, whiteList []string) {
	query := r.URL.Query()

	for k, v := range query {
		if !isWhiteListed(k, whiteList) {
			query.Del(k)
			continue
		}
		if len(v) > 1 {
			query.Set(k, v[0]) // first value
			// query.Set(k, v[len(v)-1]) // last value
		}
	}

	r.URL.RawQuery = query.Encode()
}

func filterBodyParams(r *http.Request, whiteList []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range r.Form {
		if !isWhiteListed(k, whiteList) {
			r.Form.Del(k)
			continue
		}
		if len(v) > 1 {
			r.Form.Set(k, v[0]) // first value
			// r.Form.Set(k, v[len(v)-1]) // last value
		}
	}
}

func isWhiteListed(param string, whiteList []string) bool {
	return slices.Contains(whiteList, param)
}
