package app

import (
	_ "fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.marqeta.com/ecray/avdb/app/handler"
	"github.marqeta.com/ecray/avdb/app/middleware"
	"github.marqeta.com/ecray/avdb/app/model"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize() {
	// set db string from os.Env
	db, err := gorm.Open(
		"postgres",
		os.ExpandEnv("host=${DB_HOST} port=5432 user=${DB_USER} dbname=${DB_NAME} sslmode=disable password=${DB_PASS}"))

	if err != nil {
		panic("Could not connect database")
	}

	// Generate schema
	a.DB = model.DBMigrate(db)
	model.PopulateAuth(a.DB)
	// Create new router
	a.Router = mux.NewRouter()
	// Set handlers
	a.setRouters()
}

func (a *App) setRouters() {
	// set up middleware
	//mw := ChainMiddleware(BasicAuth)
	mw := middleware.BasicAuth

	// Routing for host functions
	a.Get("/hosts", a.GetAllHosts)
	a.Get("/hosts/{name}", a.GetHost)
	a.Post("/hosts/{name}", mw(a.CreateHost, a.DB))
	a.Put("/hosts/{name}", mw(a.UpdateHost, a.DB))
	a.Delete("/hosts/{name}", mw(a.DeleteHost, a.DB))
	// Routing for group functions
	a.Get("/groups", a.GetAllGroups)
	a.Get("/groups/{name}", a.GetGroup)
	a.Post("/groups/{name}", mw(a.CreateGroup, a.DB))
	a.Put("/groups/{name}", mw(a.UpdateGroup, a.DB))
	a.Delete("/groups/{name}", mw(a.DeleteGroup, a.DB))
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

// App run
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
