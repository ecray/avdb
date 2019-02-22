package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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

func sliceContains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func debugBody(data *Request) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(data)
}

// this removes hosts when entry has -, ie -web02
func removeByDash(o []string, s string) []string {
	s = strings.TrimLeft(s, "-")

	// get index in original
	idx := findIndex(o, s)

	// delete from original
	o[len(o)-1], o[idx] = o[idx], o[len(o)-1]
	return o[:len(o)-1]
}

func findIndex(s []string, x string) int {
	for i, z := range s {
		if x == z {
			return i
		}
	}
	return len(s)
}
