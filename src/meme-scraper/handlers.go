package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getQueryParams(r *http.Request) (map[string][]string, error) {
	qp := r.URL.Query()
	ret := make(map[string][]string)
	for k, v := range qp {
		ret[k] = v
		key := k
		if strings.HasPrefix(k, "views") || strings.HasPrefix(k, "videos") ||
			strings.HasPrefix(k, "images") || strings.HasPrefix(k, "comments") {
			if strings.HasSuffix(k, "_lt") || strings.HasSuffix(k, "_gt") {
				key = k[:len(k)-3]
			}
		}
		if strings.HasPrefix(k, "created") || strings.HasPrefix(k, "lastUpdated") {
			if strings.HasSuffix(k, "_before") {
				key = k[:len(k)-len("_before")]
			}
			if strings.HasSuffix(k, "_after") {
				key = k[:len(k)-len("_after")]
			}
		}
		switch key {
		case "name", "src":
			continue
		case "views", "videos", "images", "comments":
			if _, err := strconv.Atoi(v[0]); err != nil {
				return nil, fmt.Errorf("Cannot convert %s to integer: %v", k, err)
			}
		case "created", "lastUpdated":
			t, err := time.Parse(time.RFC3339, v[0])
			if err != nil {
				return nil, fmt.Errorf("Failed to parse time string for %s: %v", k, err)
			}
			if strings.HasSuffix(k, "_before") {
				key = key + "_lt"
			}
			if strings.HasSuffix(k, "_after") {
				key = key + "_gt"
			}
			delete(ret, k)
			ret[key] = []string{strconv.FormatInt(t.Unix(), 10)}
		default:
			return nil, fmt.Errorf("Unrecongnized key: %s", k)
		}
	}
	return ret, nil
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
