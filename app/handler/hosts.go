package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ecray/avdb/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// CreateHost ...
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

// GetAllHosts ...
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

// GetHost ...
func GetHost(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	host := getHostOr404(db, name, w, r)
	if host == nil {
		return
	}
	respondJSON(w, http.StatusOK, host)
}

// DeleteHost ...
func DeleteHost(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	host := getHostOr404(db, name, w, r)
	if host == nil {
		return
	}
	if err := db.Unscoped().Delete(&host).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// UpdateHost
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
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Failed: ", err, r.Body)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// debugging body data
	/* enc := json.NewEncoder(os.Stdout).SetIndent("", "    ").Encode(data) */

	// convert original host data to map to parse
	origin := make(map[string]interface{})
	if err := json.Unmarshal(host.Data.RawMessage, &origin); err != nil {
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
	var b []byte
	b, err := json.Marshal(&origin)
	if err != nil {
		log.Printf("failed to marshal")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	host.Data.RawMessage = b

	if err := db.Save(&host).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, host)
}

func getHostOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Host {
	host := model.Host{}
	if err := db.First(&host, model.Host{Name: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &host
}
