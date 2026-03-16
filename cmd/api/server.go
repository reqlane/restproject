package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	mw "restproject/internal/api/middlewares"
	"strconv"
	"strings"
	"sync"
)

type teacher struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

var (
	teachers = make(map[int]teacher)
	mu       = &sync.Mutex{}
	nextID   = 1
)

// Initialize some data
func init() {
	teachers[nextID] = teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++
	teachers[nextID] = teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "11A",
		Subject:   "Biology",
	}
	nextID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		mu.Lock()
		teacherList := make([]teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}
		mu.Unlock()

		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// Handle Path parameter
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		mu.Lock()
		teacher, exists := teachers[id]
		mu.Unlock()
		if !exists {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(teacher)
	}
}

func postTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	addedTeachers := make([]teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Root Route"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Teachers Route"))
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Students Route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Students Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Students Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Students Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Students Route"))
	}
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET Method on Execs Route"))
	case http.MethodPost:
		w.Write([]byte("Hello POST Method on Execs Route"))
	case http.MethodPut:
		w.Write([]byte("Hello PUT Method on Execs Route"))
	case http.MethodPatch:
		w.Write([]byte("Hello PATCH Method on Execs Route"))
	case http.MethodDelete:
		w.Write([]byte("Hello DELETE Method on Execs Route"))
	}
}

func main() {

	port := 3000
	cert := "cert.pem"
	key := "key.pem"
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := mw.NewRateLimiter(5, time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// secureMux := applyMiddlewares(mux,
	// 	mw.ResponseTime,
	// 	rl.RateLimit,
	// 	mw.Cors,
	// 	mw.SecurityHeaders,
	// 	mw.Hpp(hppOptions),
	// 	mw.Compression,
	// )
	secureMux := mw.SecurityHeaders(mux)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Handler: mux,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}

type middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
