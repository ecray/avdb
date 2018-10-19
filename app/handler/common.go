package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(resp))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"Error": message})
}

func queryBuilder(strs url.Values) string {
	var sb strings.Builder
	i := 0
	for k, v := range strs {
		// if query is more than one key append, or seek operator
		if i >= 1 {
			sb.WriteString(" AND ")
		}

		// if query is hosts, check hosts column, else check data
		if k == "hosts" {
			sb.WriteString(fmt.Sprintf("'%s' = ANY(%s)", v[0], k))
		} else if k != "op" {
			sb.WriteString(fmt.Sprintf("data->>'%s' = '%s'", k, v[0]))
		}
		i++
	}
	return sb.String()
}
