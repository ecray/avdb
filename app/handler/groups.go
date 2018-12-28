package handler

import (
	"encoding/json"
	_ "fmt"
	"github.com/ecray/avdb/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Data  map[string]string `json:"data"`
	Hosts []string          `json:"hosts"`
}

func CreateGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group := model.Group{}

	group.Name = vars["name"]
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		switch {
		case err == io.EOF:
			// empty body; not an error, client didn't send data
		case err != nil:
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	defer r.Body.Close()

	// Debugging
	// debugBody(group)

	if err := db.Save(&group).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value"):
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(w, http.StatusCreated, group)
}

func GetAllGroups(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	groups := []model.Group{}
	query := r.URL.Query()
	log.Println(query)
	//set to true to log query
	db.LogMode(false)

	// if query params found, build query
	if len(query) > 0 {
		db.Where(queryBuilder(query)).Find(&groups)
	} else {
		db.Find(&groups)
	}
	respondJSON(w, http.StatusOK, groups)
}

func GetGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	group := getGroupOr404(db, name, w, r)
	if group == nil {
		return
	}
	respondJSON(w, http.StatusOK, group)
}

func DeleteGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	group := getGroupOr404(db, name, w, r)
	if group == nil {
		return
	}
	if err := db.Delete(&group).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func UpdateGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Verify group exists
	name := vars["name"]
	group := getGroupOr404(db, name, w, r)
	if group == nil {
		return
	}

	// decode response data into map
	var data *Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Failed response decode: ", err, r.Body)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// debug request body
	//debugBody(data)

	// convert original group data to map to iterate
	origin := make(map[string]interface{})
	err = json.Unmarshal(group.Data.RawMessage, &origin)
	if err != nil {
		log.Println("error", err)
	}

	/* Update origin with new data, and cull out nil kv
	   This needs to check for nested data updates, otherwise
	   it will blow away entire array/maps */
	for k, v := range data.Data {
		// Get data from response
		if v == "" {
			_, ok := origin[k]
			if ok {
				log.Println("Deleting ", k)
				delete(origin, k)
			}
		} else {
			origin[k] = v
		}
	}
	for _, v := range data.Hosts {
		// Get data from response
		if v == "" {
			log.Println("Found invalid data in request")
			break
		}
		ok := sliceContains(v, group.Hosts)
		if ok {
			//log.Println("Found existing entry ", v)
		} else if strings.HasPrefix(v, "-") {
			// Remove host from original
			group.Hosts = removeByDash(group.Hosts, v)
		} else {
			// update group
			group.Hosts = append(group.Hosts, v)
		}
	}

	// Convert back to model
	group.Data.RawMessage, err = json.Marshal(&origin)
	if err != nil {
		log.Printf("failed to marshal")
	}

	if err := db.Save(&group).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, group)
}

func getGroupOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Group {
	group := model.Group{}
	err := db.Where(model.Group{Name: name}).Find(&group).Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &group
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
