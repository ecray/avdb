package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	//db.LogMode(true)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	tag := model.Tag{
		Name: vars["name"],
		Host: data.Host,
	}

	// Check for existing tag combo
	// If match or not, save
	db.Where("tag = ? AND host = ?", tag.Name, tag.Host).Find(&tag)
	if err := db.Save(&tag).Error; err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value"):
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondJSON(w, http.StatusCreated, tag)
}

// GetAllTags ...
func GetAllTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tags := []model.Tag{}
	query := r.URL.Query()
	//db.LogMode(true)
	var results []string
	fmt.Println(query.Get("host"))
	// if host query
	if len(query) > 0 {
		db.Where("host = ?", query.Get("host")).Find(&tags)
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
	tags := []model.Tag{}
	var results []string

	// debug queries
	//db.LogMode(true)

	// Search for tag and populate results
	db.Model(&model.Tag{}).Where("tag = ?", name).Find(&tags)
	if len(tags) == 0 {
		respondError(w, http.StatusNotFound, fmt.Sprint("No Record Found For ", name))
		return
	}
	for _, t := range tags {
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

	//db.LogMode(true)

	// Unpack body to get host
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Where("tag = ? AND host = ?", name, data.Host).First(&tag).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := db.Unscoped().Delete(&tag).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, "")
}
