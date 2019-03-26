package handler

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

func GetHealth(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if _, err := db.Raw("SELECT 1;").Rows(); err != nil {
		http.Error(w, "ERROR: Database connection not available", http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("OK"))
}
