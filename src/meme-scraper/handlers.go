package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getQueryParams(r *http.Request) (map[string][]string, error) {
	qp := r.URL.Query()
	for k, v := range qp {
		key := k
		if strings.HasSuffix(k, "_lt") || strings.HasSuffix(k, "_gt") {
			key = k[:len(k)-3]
		}
		switch key {
		case "name", "src":
			continue
		case "views", "videos", "images", "comments":
			if _, err := strconv.Atoi(v[0]); err != nil {
				return nil, fmt.Errorf("Cannot convert %s to integer: %v", k, err)
			}
		case "created", "lastUpdated":
			// TODO check time
		default:
			return nil, fmt.Errorf("Unrecongnized key: %s", k)
		}
	}
	return qp, nil
}

func memeHandler(w http.ResponseWriter, r *http.Request) {
	queryParams, err := getQueryParams(r)
	if err != nil {
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
