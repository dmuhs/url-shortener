package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	config "./config"
	helpers "./helpers"
	"github.com/gorilla/mux"
)

func getRequestContext(r *http.Request) config.Configuration {
	// Value() always returns an interface{} type, so we have to cast it
	// to our Configuration struct again to make it usable
	return r.Context().Value(contextKey("Configuration")).(config.Configuration)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	configCtx := getRequestContext(r)

	var sb helpers.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&sb)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := helpers.CheckAndGetShortCode(sb.URL, configCtx.SQLitePath, configCtx.ShortCodeLength)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// construct the JSON response with long and short URL
	resp := helpers.ConstructRedirectResponse(sb.URL, configCtx.Domain, configCtx.Port, code)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func redirectToURL(w http.ResponseWriter, r *http.Request) {
	configCtx := getRequestContext(r)
	code := mux.Vars(r)["code"]

	url, err := helpers.GetLongURLByShortCode(configCtx.SQLitePath, code)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func displayDefaultText(w http.ResponseWriter, r *http.Request) {
	configCtx := getRequestContext(r)
	database, err := helpers.GetDatabaseConnector(configCtx.SQLitePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var rowCount int
	row := database.QueryRow("SELECT COUNT(*) FROM urls")
	err = row.Scan(&rowCount)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "There are %d rows in the database", rowCount)
}
