package middleware

import (
	_ "fmt"
	"log"
	"net/http"
	_ "os"

	"github.com/jinzhu/gorm"

	"github.marqeta.com/ecray/avdb/app/model"
)

//type middleware func(http.HandlerFunc) http.HandlerFunc

func BasicAuth(next http.HandlerFunc, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("Authenticating request to %s", r.RequestURI)
		// get header/token, check model/db and validate
		auth := model.Auth{}
		token := r.Header.Get("Auth-Token")
		err := db.First(&auth, model.Auth{Token: token}).Error
		if auth.Token == token {
			next.ServeHTTP(w, r)
		} else if err != nil {
			http.Error(w, "Error", http.StatusForbidden)
		} else {
			log.Printf("Authentication failed for request to %s", r.RequestURI)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}
