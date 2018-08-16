package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tkanos/gonfig"

	config "./config"
	helpers "./helpers"
)

func loadConfiguration(path string) config.Configuration {
	configuration := config.Configuration{}
	err := gonfig.GetConf("config.json", &configuration)
	if err != nil {
		log.Fatal(err)
	}
	return configuration
}

// use a separate context Key type to avoid collisions across packages, as
// pkg1.contextKey("foo") != pkg2.contextKey("foo")
type contextKey string

// ContextMiddleware to add the configuration struct to all incoming requests
func ContextMiddleware(next http.Handler, c config.Configuration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKey("Configuration"), c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	configuration := loadConfiguration("config.json")
	database, err := helpers.GetDatabaseConnector(configuration.SQLitePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing table 'urls'")
	statement, _ := database.Prepare(`
    CREATE TABLE IF NOT EXISTS urls (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      long TEXT,
      short TEXT,
      added TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
	statement.Exec()
	log.Println("Setting 'urls' indices")
	statement, _ = database.Prepare("CREATE UNIQUE INDEX IF NOT EXISTS idx_long ON urls (long)")
	statement.Exec()
	database.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", displayDefaultText).Methods("GET")
	router.HandleFunc("/", shortenURL).Methods("POST")
	router.HandleFunc("/{code}", redirectToURL).Methods("GET")
	log.Printf("Listening on port %d", configuration.Port)
	var host string
	if configuration.Public {
		host = ":" + strconv.Itoa(configuration.Port)
	} else {
		host = "localhost:" + strconv.Itoa(configuration.Port)
	}
	// add configuration context middleware to router
	contextRouter := ContextMiddleware(router, configuration)
	// crash loudly during deployment if something goes wrong
	log.Fatal(http.ListenAndServe(host, contextRouter))
}
