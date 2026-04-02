package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// URL Path
		sanitizedPath, err := sanitize(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Query params
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for key, values := range params {
			sanitizedKey, err := sanitize(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var sanitizedValues []string
			for _, value := range values {
				sanitizedValue, err := sanitize(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, sanitizedValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
		}
		r.URL.Path = sanitizedPath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		// Body
		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "error reading request body", http.StatusBadRequest)
					return
				}
				bodyString := strings.TrimSpace(string(bodyBytes))

				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))

				if len(bodyString) > 0 {
					var inputData any
					err := json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						http.Error(w, "invalid JSON body", http.StatusBadRequest)
						return
					}

					sanitizedData, err := sanitize(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					sanitizedBody, err := json.Marshal(sanitizedData)
					if err != nil {
						http.Error(w, "error sanitizing body", http.StatusBadRequest)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))
					r.ContentLength = int64(len(sanitizedBody))
				}
			}
		} else if r.Header.Get("Content-Type") != "" {
			log.Printf("Received request with unsuported Content-Type: %s. Expected application/json.\n", r.Header.Get("Content-Type"))
			http.Error(w, "unsuported Content-Type, please use application/json.", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func sanitize(data any) (any, error) {
	switch v := data.(type) {
	case string:
		return sanitizeString(v), nil
	case []any:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v, nil
	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", data)
	}
}

func sanitizeValue(value any) any {
	switch v := value.(type) {
	case string:
		return sanitizeString(v)
	case []any:
		for i, value := range v {
			v[i] = sanitizeValue(value)
		}
		return v
	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v
	default:
		return v
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
