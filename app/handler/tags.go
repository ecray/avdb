package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecray/avdb/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// CreateTag ...
func CreateTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//tag := model.Tag{}
	var data struct {
		Host string `json:"host"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// Fetch host record
	host := getHostID(db, data.Host, w, r)

	tag := model.Tag{
		Name:   vars["name"],
		HostID: host.ID,
	}

	// Check for existing tag combo
	// If match or not, save
	db.Where("tag = ? AND host_id = ?", tag.Name, tag.HostID).Find(&tag)
	if err := db.Save(&tag).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, tag)
}

// GetAllTags ...
func GetAllTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tags := []model.Tag{}
	query := r.URL.Query()
	var results []string

	if len(query) > 0 {
		// if host query
		name := query.Get("host")
		if name != "" {
			host := getHostID(db, name, w, r)
			db.Where("host_id = ?", host.ID).Find(&tags)
		} else {
			respondError(w, http.StatusNotImplemented, "Only host query supported")
			return
		}
	} else {
		db.Select("DISTINCT(tag)").Find(&tags)
	}

	for _, t := range tags {
		results = append(results, t.Name)
	}
	respondJSON(w, http.StatusOK, results)
}

// GetTag ...
func GetTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	type temp struct {
		Host string
	}
	var temps []temp
	var results []string

	// Search for tag and populate results
	// GORM bug with Raw, Rows keeps giving 0 results
	// short term - move to plain SQL
	db.Raw("SELECT hosts.host FROM tags INNER JOIN hosts ON tags.host_id = hosts.id WHERE tag = ?", name).Scan(&temps)
	if len(temps) == 0 {
		respondError(w, http.StatusNotFound, fmt.Sprint("No Record Found For ", name))
		return
	}
	for _, t := range temps {
		results = append(results, t.Host)
	}
	respondJSON(w, http.StatusOK, results)
}

// DeleteTag ...
func DeleteTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := model.Tag{}
	name := vars["name"]

	var data struct {
		Host string `json:"host"`
	}

	// Unpack body to get host
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	// Fetch host record
	host := getHostID(db, data.Host, w, r)

	// Fetch record
	if err := db.Where("tag = ? AND host_id = ?", name, host.ID).First(&tag).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	// Delete record
	if err := db.Unscoped().Delete(&tag).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, "")
}

func getHostID(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *model.Host {
	host := model.Host{}
	if err := db.First(&host, model.Host{Name: name}).Error; err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprint("No Record Found For ", name))
		return nil
	}
	return &host
}
