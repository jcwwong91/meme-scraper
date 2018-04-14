package main

import (
	"fmt"
	"strings"
	"time"
)

func genQuery(qp map[string][]string) string {
	queryString := "SELECT * FROM memeInfo"
	if len(qp) != 0 {
		queryString += " WHERE "
		first := true
		for k, v := range qp {
			if !first {
				queryString += " AND "
			}
			// Parser seems to always leave a single value for v
			val := strings.Replace(v[0], "'", "''", -1)
			if strings.HasSuffix(k, "_lt") {
				queryString += k[:len(k)-3] + " < '" + val + "'"
			} else if strings.HasSuffix(k, "_gt") {
				queryString += k[:len(k)-3] + " > '" + val + "'"
			} else {
				if k == "name" || k == "src" {
					queryString += k + " LIKE '%" + val + "%'"
				} else {
					queryString += k + " = '" + val + "'"
				}
			}
			first = false
		}
	}
	fmt.Println(queryString)
	return queryString
}

func query(qp map[string][]string) ([]meme, error) {
	query := genQuery(qp)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Query Failed: %v", err)
	}
	defer rows.Close()

	ret := []meme{}
	cols, _ := rows.Columns()
	fmt.Println(cols)
	for rows.Next() {
		var m meme
		var created, lastUpdated int64

		err = rows.Scan(&m.Name, &m.Src, &m.Views, &m.Videos, &m.Images, &m.Comments, &created, &lastUpdated)
		if err != nil {
			return nil, fmt.Errorf("Failed to load meme struct: %v", err)
		}
		m.Created = time.Unix(created, 0)
		m.LastUpdated = time.Unix(lastUpdated, 0)

		ret = append(ret, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Failed to load meme struct: %v", err)
	}
	return ret, nil
}
