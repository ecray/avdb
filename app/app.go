package app

import (
	"log"
	"net/http"
	"os"
	"time"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/ecray/avdb/app/handler"
	"github.com/ecray/avdb/app/middleware"
	"github.com/ecray/avdb/app/model"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Run(host string) {
	var db *gorm.DB
	var err error
	var dur int = 5

	// Set up db connection with retries
	connString := os.ExpandEnv("host=${DB_HOST} port=${DB_PORT} user=${DB_USER} dbname=${DB_NAME} sslmode=disable password=${DB_PASS}")
	// DB connectivity check every 5 seconds, total 25 seconds
	for i := 5; i > 0; i-- {
		db, err = gorm.Open("postgres", connString)
		if err != nil && i == 0 {
			panic("Connection failed after retries, abandoning.")
		} else if err != nil {
			log.Printf("Connection Failed...Retries %d left - sleeping %d seconds..", i, dur)
			time.Sleep(time.Duration(dur) * time.Second)
		}
	}
	defer db.Close()

	// Generate schema
	log.Println("Auto-migrating schema...")
	a.DB = model.DBMigrate(db)

	// Populate initial token and log
	log.Println("Checking token credential..")
	model.PopulateAuth(a.DB)

	// Create new router
	a.Router = mux.NewRouter()

	// Set handlers
	a.setRouters()

	// Log and Serve
	logged := gh.LoggingHandler(os.Stdout, a.Router)
	http.ListenAndServe(host, logged)
}

func (a *App) setRouters() {
	// set up auth middleware
	mwauth := middleware.BasicAuth

	// Routing for host functions
	a.Get("/api/v1/hosts", a.GetAllHosts)
	a.Get("/api/v1/hosts/{name}", a.GetHost)
	a.Post("/api/v1/hosts/{name}", mwauth(a.CreateHost, a.DB))
	a.Put("/api/v1/hosts/{name}", mwauth(a.UpdateHost, a.DB))
	a.Delete("/api/v1/hosts/{name}", mwauth(a.DeleteHost, a.DB))
	// Routing for group functions
	a.Get("/api/v1/groups", a.GetAllGroups)
	a.Get("/api/v1/groups/{name}", a.GetGroup)
	a.Post("/api/v1/groups/{name}", mwauth(a.CreateGroup, a.DB))
	a.Put("/api/v1/groups/{name}", mwauth(a.UpdateGroup, a.DB))
	a.Delete("/api/v1/groups/{name}", mwauth(a.DeleteGroup, a.DB))
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Host handlers
func (a *App) GetAllHosts(w http.ResponseWriter, r *http.Request) {
	handler.GetAllHosts(a.DB, w, r)
}

func (a *App) CreateHost(w http.ResponseWriter, r *http.Request) {
	handler.CreateHost(a.DB, w, r)
}

func (a *App) GetHost(w http.ResponseWriter, r *http.Request) {
	handler.GetHost(a.DB, w, r)
}

func (a *App) DeleteHost(w http.ResponseWriter, r *http.Request) {
	handler.DeleteHost(a.DB, w, r)
}

func (a *App) UpdateHost(w http.ResponseWriter, r *http.Request) {
	handler.UpdateHost(a.DB, w, r)
}

// Group handlers
func (a *App) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	handler.GetAllGroups(a.DB, w, r)
}

func (a *App) CreateGroup(w http.ResponseWriter, r *http.Request) {
	handler.CreateGroup(a.DB, w, r)
}

func (a *App) GetGroup(w http.ResponseWriter, r *http.Request) {
	handler.GetGroup(a.DB, w, r)
}

func (a *App) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	handler.DeleteGroup(a.DB, w, r)
}

func (a *App) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	handler.UpdateGroup(a.DB, w, r)
}
