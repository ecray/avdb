package handler

import (
	"encoding/json"
	_ "fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.marqeta.com/ecray/avdb/app/model"
	"io"
	"log"
	"net/http"
	_ "os"
	"strings"
)

func CreateHost(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	host := model.Host{}

	host.Name = vars["name"]
	if err := json.NewDecoder(r.Body).Decode(&host.Data); err != nil {
		switch {
		case err == io.EOF:
			// empty body; client didn't send data
		case err != nil:
			// log.Println("Failed to decode data: ", err, r.Body)
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	defer r.Body.Close()

	// Debugging
	// debugBody(host)

	if err := db.Save(&host).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value"):
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(w, http.StatusCreated, host)
}

func GetAllHosts(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	hosts := []model.Host{}
	query := r.URL.Query()
	db.LogMode(false)

	// if query params found, build query
	if len(query) > 0 {
		db.Where(queryBuilder(query)).Find(&hosts)
	} else {
		db.Find(&hosts)
	}
	respondJSON(w, http.StatusOK, hosts)
}

func GetHost(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	host := getHostOr404(db, name, w, r)
	if host == nil {
		return
	}
	respondJSON(w, http.StatusOK, host)
}

func DeleteHost(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	host := getHostOr404(db, name, w, r)
	if host == nil {
		return
	}
	if err := db.Delete(&host).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func UpdateHost(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Verify host exists
	name := vars["name"]
	host := getHostOr404(db, name, w, r)
	if host == nil {
		return
	}

	// decode response data into map
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("Failed: ", err, r.Body)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	/* enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(data) */

	// convert original host data to map to parse
	origin := make(map[string]interface{})
	err = json.Unmarshal(host.Data.RawMessage, &origin)
	if err != nil {
		log.Println("error", err)
	}
	/* Update origin with new data, and cull out nil kv
	   This needs to check for nested data updates, otherwise
	   it will blow away entire array/maps
	*/
	for k, v := range data {
		if v == nil {
			_, ok := origin[k]
			if ok {
				log.Println("Deleting ", k)
				delete(origin, k)
			}
		} else {
			origin[k] = v
		}
	}

	// Convert back to model
	host.Data.RawMessage, err = json.Marshal(&origin)
	if err != nil {
		log.Printf("failed to marshal")
	}

	if err := db.Save(&host).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, host)
}

func getHostOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Host {
	host := model.Host{}
	err := db.First(&host, model.Host{Name: name}).Error
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &host
}
