package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func validateQueryParams(qp map[string][]string) error {
	// TODO
	return nil
}

func memeHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	if err := validateQueryParams(queryParams); err != nil {
		fmt.Fprintf(w, "Bad Parameters: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	results, err := query(queryParams)
	if err != nil {
		fmt.Fprintf(w, "Failed to get results: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintf(w, "Failed to convert data to JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", string(b))
}
