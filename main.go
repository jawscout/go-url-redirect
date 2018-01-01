package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	staticPages = map[string]func(http.ResponseWriter, *http.Request) error{
		"icon.ico":   nil,
		"robots.txt": nil,
	}
	// temporaryRedirects This will eventually be replased by Google cloud sql (or the DB option that is included in the free teir)
	temporaryRedirects = map[string]string{
		"github": "https://github.com/jawscout",
		"home":   "https://www.jawscout.cc",
	}
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{key}", keyHandler)
	srv := &http.Server{
		Handler: r,
	}
	log.Fatal(srv.ListenAndServe())
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value, err := lookupKey(vars["key"])
	if handleError(w, r, err) {
		return
	}
	// Check to see if key is a static value
	if handler, found := staticPages[value]; found {
		// w.Write([]byte("Static resouse found"))
		if handler != nil {
			handler(w, r)
		}
		return
	}
	http.Redirect(w, r, value, http.StatusFound)
}

func lookupKey(key string) (string, error) {
	if val, ok := temporaryRedirects[key]; ok {
		return val, nil
	}
	return "", errors.New("Key not in DB")
}

func handleError(w http.ResponseWriter, r *http.Request, e error) bool {
	// If error is nil stop here
	if e == nil {
		return false
	}
	switch e := e.(type) {
	case error:
		// We can retrieve the status here and write out a specific
		// HTTP status code.
		log.Printf("HTTP %d - %s", http.StatusInternalServerError, e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
	default:
		// Any error types we don't specifically look out for default
		// to serving a HTTP 500
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return true
}
