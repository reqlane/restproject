package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type HPPConfig struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	WhiteList                   []string
}

func Hpp(config HPPConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.CheckQuery && r.URL.Query() != nil {
				filterQueryParams(r, config.WhiteList)
			}

			if config.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, config.CheckBodyOnlyForContentType) {
				filterBodyParams(r, config.WhiteList)
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
			query.Set(k, v[0])
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
			r.Form.Set(k, v[0])
		}
	}
}

func isWhiteListed(param string, whiteList []string) bool {
	return slices.Contains(whiteList, param)
}
